package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"
	"user-service/models/dto"
	"user-service/pkg/base"
	"user-service/pkg/config"
	"user-service/pkg/errors"
	"user-service/repositories"

	"user-service/pkg/logger"
	"user-service/pkg/redis"

	apperrors "user-service/pkg/errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	Login(ctx context.Context, req *dto.LoginInput) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	HashPassword(password string) (string, error)

	// OTP methods
	GenerateAndSendOTP(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, email, otp string) error
	ResendOTP(ctx context.Context, email string) error
}

type AuthService struct {
	*base.BaseService
	userRepo     repositories.IUserRepository
	redisCli    *redis.Client
	emailService IMailtrapEmailService
	memberService IMemberService
	instructorService IInstructorService
	roleService   IRoleService
}

func NewAuthService(userRepo repositories.IUserRepository, redisCli *redis.Client, emailService IMailtrapEmailService, memberService IMemberService, instructorService IInstructorService, roleService IRoleService) IAuthService {
	return &AuthService{
		userRepo:        userRepo,
		redisCli:        redisCli,
		emailService:    emailService,
		memberService:   memberService,
		instructorService: instructorService,
		roleService:     roleService,
	}
}

type Claims struct {
	User *dto.GetUserResponse
	jwt.RegisteredClaims
}

// Login authenticates a user with email and password
func (s *AuthService) Login(ctx context.Context, req *dto.LoginInput) (*dto.LoginResponse, error) {
	// Find user by email or username
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		// If not found by email, try by username
		user, err = s.userRepo.FindByUsername(ctx, req.Email)
		if err != nil {
			return nil, apperrors.ErrInvalidCredentials
		}
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	// Check if user is verified (optional - uncomment if required)
	// if !user.IsVerified {
	// 	return nil, apperrors.ErrEmailNotVerified
	// }

	cfg := config.Get()

	// s.LogError("error happening auth service:", logger.LogField("exipry hour", cfg.JWT.ExpiryHour));
	// s.LogError("error happening auth service:", logger.LogField("exipry days", cfg.JWT.RefreshTokenExpiryDays));
	// Set expiration time
	expirationTime := time.Now().Add(time.Duration(cfg.JWT.ExpiryHour) * time.Hour).Unix()

	// Build user response
	userResp := dto.GetUserResponse{
		UserID:      user.ID,
		Email:       user.EmailAddress,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Image:       user.Image,
		DateOfBirth: user.DateOfBirth,
		Address:     user.Address,
		RoleID:      user.RoleID,
	}

	// Create JWT claims
	claims := &Claims{
		User: &userResp,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
		},
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	// Generate refresh token
	refreshToken := generateRefreshToken()
	refreshKey := fmt.Sprintf("refresh_token:%s", refreshToken)
	expiryDuration := time.Duration(cfg.JWT.RefreshTokenExpiryDays) * 24 * time.Hour
	if err := s.redisCli.Client.Set(ctx, refreshKey, user.ID.String(), expiryDuration).Err(); err != nil {
		return nil, apperrors.ErrInternalServer
	}

	return &dto.LoginResponse{
		User:         userResp,
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    expirationTime,
	}, nil
}

// RefreshToken validates refresh token and returns new access + refresh tokens
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	cfg := config.Get()

	// Get stored refresh token from Redis
	refreshKey := fmt.Sprintf("refresh_token:%s", refreshToken)
	storedUserID, err := s.redisCli.Client.Get(ctx, refreshKey).Result()
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	// Parse stored user ID
	userID, err := uuid.Parse(storedUserID)
	if err != nil {
		return nil, apperrors.ErrUnauthorized
	}

	// Find user by ID
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, apperrors.ErrUserNotFound
	}

	// Generate new access token
	accessExpirationTime := time.Now().Add(time.Duration(cfg.JWT.ExpiryHour) * 3600).Unix()
	userResp := dto.GetUserResponse{
		UserID:      user.ID,
		Email:       user.EmailAddress,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Image:       user.Image,
		DateOfBirth: user.DateOfBirth,
		Address:     user.Address,
		RoleID:      user.RoleID,
	}

	claims := &Claims{
		User: &userResp,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(accessExpirationTime, 0)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}

	// Delete old refresh token
	s.redisCli.Client.Del(ctx, refreshKey)

	// Generate new refresh token
	newRefreshToken := generateRefreshToken()
	refreshKey = fmt.Sprintf("refresh_token:%s", newRefreshToken)
	expiryDuration := time.Duration(cfg.JWT.RefreshTokenExpiryDays) * 24 * time.Hour
	if err := s.redisCli.Client.Set(ctx, refreshKey, user.ID.String(), expiryDuration).Err(); err != nil {
		return nil, apperrors.ErrInternalServer
	}

	return &dto.LoginResponse{
		User:         userResp,
		AccessToken:  accessTokenString,
		RefreshToken: newRefreshToken,
		ExpiresIn:    accessExpirationTime,
	}, nil
}

// generateRefreshToken generates a cryptographically secure random token
func generateRefreshToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	s.LogInfo("Register: Starting registration for email: " + req.Email)

	if s.isUsernameExist(ctx, req.Username) {
		s.LogInfo("Register: Username exists: " + req.Username)
		return nil, apperrors.ErrUsernameExist
	}
	s.LogInfo("Register: Username check passed")

	if s.isEmailExist(ctx, req.Email) {
		s.LogInfo("Register: Email exists: " + req.Email)
		return nil, apperrors.ErrEmailExist
	}
	s.LogInfo("Register: Email check passed")

	if req.Password != req.ConfirmPassword {
		return nil, apperrors.ErrPasswordDoesNotMatch
	}

	s.LogInfo("Register: Hashing password...")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	confirmHashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	s.LogInfo("Register: Password hashed")

	registerReq := &dto.RegisterRequest{
		FirstName:         req.FirstName,
		LastName:          req.LastName,
		Username:     req.Username,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		DateOfBirth:  req.DateOfBirth,
		Password:     string(hashedPassword),
		ConfirmPassword: string(confirmHashedPassword),
		RoleID:       req.RoleID,
	}

	s.LogInfo("Register: Fetching roles...")
	roles, err := s.roleService.FindAllRoles(ctx)
	if err != nil {
		s.LogError("error happening role repo:", logger.LogField("error", err));
		return nil, err
	}
	s.LogInfo("Register: Roles fetched, count: " + fmt.Sprintf("%d", len(roles)))

	s.LogInfo("Register: Creating user in database...")
	user, err := s.userRepo.Create(ctx, registerReq)
	if err != nil {
		s.LogError("error happening user repo:", logger.LogField("error", err));
		return nil, err
	}
	s.LogInfo("Register: User created, ID: " + user.ID.String())

	s.LogInfo("Register: Creating profile for role ID: " + fmt.Sprintf("%d", req.RoleID))
	for _, role := range roles {
		if role.ID == req.RoleID {
			s.LogInfo("Register: Matched role: " + role.Name)
			if(strings.ToLower(role.Name) == "member") {
				s.LogInfo("Register: Creating member profile...")
				if _, err := s.memberService.CreateMemberProfile(ctx, user.ID); err != nil {
					return nil, fmt.Errorf("failed to create member profile: %w", err)
				}
				s.LogInfo("Register: Member profile created")
			}else if(strings.ToLower(role.Name) == "instructor") {
				s.LogInfo("Register: Creating instructor profile...")
				if _, err := s.instructorService.CreateInstructorProfile(ctx, user.ID); err != nil {
					return nil, fmt.Errorf("failed to create instructor profile: %w", err)
				}
				s.LogInfo("Register: Instructor profile created")
			}
			break
		}
	}
	s.LogInfo("Register: Profile creation completed")

	s.LogInfo("Register: Starting async OTP generation and sending...")
	go func() {
		s.LogInfo("Register: Inside goroutine - generating OTP for: " + user.EmailAddress)
		if err := s.GenerateAndSendOTP(context.Background(), user.EmailAddress); err != nil {
			s.LogError("Failed to send OTP after registration", logger.LogField("error", err))
		} else {
			s.LogInfo("Register: OTP sent successfully")
		}
	}()

	s.LogInfo("Register: Generating JWT token...")
	cfg := config.Get()
	accessExpirationTime := time.Now().Add(time.Duration(cfg.JWT.ExpiryHour) * time.Hour).Unix()
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(accessExpirationTime, 0)),
		},
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return nil, apperrors.ErrInternalServer
	}
	s.LogInfo("Register: JWT token generated")

	s.LogInfo("Register: Storing refresh token in Redis...")
	// Generate refresh token
	refreshToken := generateRefreshToken()
	refreshKey := fmt.Sprintf("refresh_token:%s", refreshToken)
	expiryDuration := time.Duration(cfg.JWT.RefreshTokenExpiryDays) * 24 * time.Hour
	if err := s.redisCli.Client.Set(ctx, refreshKey, user.ID.String(), expiryDuration).Err(); err != nil {
		return nil, apperrors.ErrInternalServer
	}
	s.LogInfo("Register: Refresh token stored")

	s.LogInfo("Register: Returning response...")
	return &dto.RegisterResponse{
		User: dto.CreateUserResponse{
			UserID:      user.ID,
			Email:       user.EmailAddress,
			Username:    user.Username,
			PhoneNumber: user.PhoneNumber,
			RoleID:      user.RoleID,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			DateOfBirth: user.DateOfBirth.Format("2006-01-02"),
		},
		AccessToken:  tokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    accessExpirationTime,
	}, nil
}

// ChangePassword changes a user's password
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

// HashPassword hashes a password and returns the hash
func (s *AuthService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// GenerateAndSendOTP generates a 6-digit OTP, stores it in Redis, and sends it via email
func (s *AuthService) GenerateAndSendOTP(ctx context.Context, email string) error {
	s.LogInfo("GenerateAndSendOTP: Starting for email: " + email)
	// Verify email exists
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.LogError("GenerateAndSendOTP: User not found", logger.LogField("error", err))
		return apperrors.ErrUserNotFound
	}
	s.LogInfo("GenerateAndSendOTP: User found, ID: " + user.ID.String())

	// Check if already verified
	if user.IsVerified {
		s.LogInfo("GenerateAndSendOTP: User already verified")
		return apperrors.ErrAlreadyVerified
	}

	// Generate 6-digit OTP
	s.LogInfo("GenerateAndSendOTP: Generating OTP...")
	otp, err := s.generateOTP()
	if err != nil {
		s.LogError("GenerateAndSendOTP: OTP generation failed", logger.LogField("error", err))
		return err
	}
	s.LogInfo("GenerateAndSendOTP: OTP generated")

	// Store OTP in Redis for 15 minutes with user ID key
	otpKey := fmt.Sprintf("email:otp:%s", user.ID.String())
	s.LogInfo("GenerateAndSendOTP: Storing OTP in Redis, key: " + otpKey)
	if err := s.redisCli.Client.Set(ctx, otpKey, otp, 15*time.Minute).Err(); err != nil {
		s.LogError("GenerateAndSendOTP: Failed to store OTP in redis", logger.LogField("error", err))
		return errors.ErrInternalServer
	}
	s.LogInfo("GenerateAndSendOTP: OTP stored in Redis")

	// Send OTP via email
	s.LogInfo("GenerateAndSendOTP: Sending OTP email...")
	if err := s.emailService.SendOTPEmail(ctx, email, otp); err != nil {
		s.LogError("GenerateAndSendOTP: Failed to send OTP email", logger.LogField("error", err))
		return errors.ErrInternalServer
	}
	s.LogInfo("GenerateAndSendOTP: OTP email sent successfully")

	return nil
}

// VerifyOTP verifies the OTP and marks the user as verified
func (s *AuthService) VerifyOTP(ctx context.Context, email, otp string) error {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return apperrors.ErrUserNotFound
	}

	// Check if already verified
	if user.IsVerified {
		return apperrors.ErrAlreadyVerified
	}

	// Get stored OTP from Redis
	otpKey := fmt.Sprintf("email:otp:%s", user.ID.String())
	storedOtp, err := s.redisCli.Client.Get(ctx, otpKey).Result()
	if err != nil {
		return apperrors.ErrOTPExpired
	}

	// Verify OTP matches
	if storedOtp != otp {
		return apperrors.ErrInvalidOTP
	}

	// Mark user as verified
	user.IsVerified = true
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Delete OTP from Redis after successful verification
	s.redisCli.Client.Del(ctx, otpKey)

	return nil
}

// ResendOTP generates a new OTP and sends it via email
func (s *AuthService) ResendOTP(ctx context.Context, email string) error {
	// Verify email exists
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return apperrors.ErrUserNotFound
	}

	// Check if already verified
	if user.IsVerified {
		return apperrors.ErrAlreadyVerified
	}

	// Delete existing OTP if any
	otpKey := fmt.Sprintf("email:otp:%s", user.ID.String())
	s.redisCli.Client.Del(ctx, otpKey)

	// Generate new OTP
	otp, err := s.generateOTP()
	if err != nil {
		return err
	}

	// Store new OTP in Redis for 15 minutes
	if err := s.redisCli.Client.Set(ctx, otpKey, otp, 15*time.Minute).Err(); err != nil {
		s.LogError("Failed to store OTP in redis", logger.LogField("error", err))
		return errors.ErrInternalServer
	}

	// Send new OTP via email
	if err := s.emailService.SendOTPEmail(ctx, email, otp); err != nil {
		s.LogError("Failed to send OTP email", logger.LogField("error", err))
		return errors.ErrInternalServer
	}

	return nil
}

// generateOTP generates a cryptographically secure 6-digit OTP
func (s *AuthService) generateOTP() (string, error) {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", errors.ErrGenerateOTP
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func (s *AuthService) isUsernameExist(ctx context.Context, username string) bool {
	s.LogInfo("Register: Checking if username exists: " + username)
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		s.LogInfo("Register: Username check - not found (error: " + err.Error() + ")")
		return false
	}

	if user != nil {
		s.LogInfo("Register: Username found in database")
		return true
	}

	s.LogInfo("Register: Username not found")
	return false
}

func (s *AuthService) isEmailExist(ctx context.Context, email string) bool {
	s.LogInfo("Register: Checking if email exists: " + email)
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.LogInfo("Register: Email check - not found (error: " + err.Error() + ")")
		return false
	}

	if user != nil {
		s.LogInfo("Register: Email found in database")
		return true
	}

	s.LogInfo("Register: Email not found")
	return false
}
