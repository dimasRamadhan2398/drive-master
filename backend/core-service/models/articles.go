package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// JSONB is a custom type for PostgreSQL JSONB columns
type JSONB json.RawMessage

// Value implements driver.Valuer for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

// Scan implements sql.Scanner for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	*j = bytes
	return nil
}                                                                                                                                                                                                                 
                                                                                                                                                                                                                    
  type ArticleStatus string                                                                                                                                                                                         
                                                                                                                                                                                                                    
  const (                                                                                                                                                                                                           
      ArticleStatusDraft     ArticleStatus = "draft"                                                                                                                                                                
      ArticleStatusPublished ArticleStatus = "published"                                                                                                                                                            
      ArticleStatusArchived  ArticleStatus = "archived"                                                                                                                                                             
  )                                                                                                                                                                                                                 
                                                                                                                                                                                                                    
  type Article struct {                                                                                                                                                                                             
      ID          uuid.UUID      `json:"id"           gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`                                                                                                        
                                                                                                                                                                                                                    
      // === Content Fields ===                                                                                                                                                                                     
      Title       string         `json:"title"        gorm:"size:255;not null"`                                                                                                                                     
      Slug        string         `json:"slug"          gorm:"size:255;uniqueIndex;not null"`                                                                                                                        
      LeadParagraph string       `json:"leadParagraph" gorm:"type:text"`         // highlighted paragraph                                                                                                           
      BodyBlocks  datatypes.JSON      `json:"bodyBlocks"      gorm:"type:jsonb"`
      Footer      string         `json:"footer"        gorm:"type:text"`         // author info, related links                                                                                                      
                                                                                                                                                                                                                    
      // === SEO Fields ===                                                                                                                                                                                         
      MetaTitle       string    `json:"metaTitle"     gorm:"size:70"`             // max 70 chars                                                                                                                   
      MetaDescription string    `json:"metaDescription" gorm:"size:160"`         // max 160 chars                                                                                                                   
      MetaKeywords    string    `json:"metaKeywords"  gorm:"size:500"`                                                                                                                                              
      CanonicalURL    string    `json:"canonicalUrl"  gorm:"size:500"`                                                                                                                                              
                                                                                                                                                                                                                    
      // === Open Graph (Social) ===                                                                                                                                                                                
      OGTitle       string      `json:"ogTitle"       gorm:"size:95"`             // max 95 chars                                                                                                                   
      OGDescription string      `json:"ogDescription" gorm:"size:200"`             // max 200 chars                                                                                                                 
      OGImage       string      `json:"ogImage"       gorm:"size:500"`            // URL to image                                                                                                                   
                                                                                                                                                                                                                    
      // === Media ===                                                                                                                                                                                              
      FeaturedImage string      `json:"featuredImage" gorm:"size:500"`                                                                                                                                              
      ThumbnailURL  string      `json:"thumbnailUrl"  gorm:"size:500"`                                                                                                                                              
                                                                                                                                                                                                                    
      // === Categorization ===                                                                                                                                                                                     
      CategoryID   *uuid.UUID   `json:"categoryId"    gorm:"type:uuid"`                                                                                                                                             
      Category      *Category    `json:"category,omitempty" gorm:"foreignKey:CategoryID"`                                                                                                                           
      Tags          pq.StringArray `json:"tags"        gorm:"type:text[]"`        // PostgreSQL array                                                                                                               
      AuthorID      uuid.UUID    `json:"authorId"      gorm:"type:uuid;not null"`                                                                                                                                   
                                                                                                                                                                                                                    
      // === Status & Publishing ===                                                                                                                                                                                
      Status        ArticleStatus `json:"status"      gorm:"size:20;default:'draft'"`                                                                                                                               
      PublishedAt   *time.Time    `json:"publishedAt"`                                                                                                                                                              
      ScheduledAt   *time.Time    `json:"scheduledAt"`                                                                                                                                                              
                                                                                                                                                                                                                    
      // === Metrics & Attractiveness ===                                                                                                                                                                           
      ViewCount     int64        `json:"viewCount"     gorm:"default:0"`                                                                                                                                            
      LikeCount     int64        `json:"likeCount"     gorm:"default:0"`                                                                                                                                            
      ShareCount    int64        `json:"shareCount"    gorm:"default:0"`                                                                                                                                            
      ReadingTime   int          `json:"readingTime"   gorm:"default:0"`          // in minutes                                                                                                                     
                                                                                                                                                                                                                    
      IsFeatured    bool         `json:"isFeatured"    gorm:"default:false"`                                                                                                                                        
      IsSpotlight   bool         `json:"isSpotlight"   gorm:"default:false"`                                                                                                                                        
      Priority      int          `json:"priority"      gorm:"default:0"`           // higher = more prominent                                                                                                       
                                                                                                                                                                                                                    
      // === Timestamps ===                                                                                                                                                                                         
      CreatedAt     time.Time    `json:"createdAt"`                                                                                                                                                                 
      UpdatedAt     time.Time    `json:"updatedAt"`                                                                                                                                                                 
      DeletedAt     gorm.DeletedAt `json:"deletedAt"   gorm:"index"`      
      
      Language      string    `json:"language" gorm:"size:10;default:'en'"`                                                                                                                                         
                                                                                                                                                                                                                    
      // Versioning / CMS                                                                                                                                                                                           
      Version       int       `json:"version" gorm:"default:1"`                                                                                                                                                     
                                                                                                                                                                                                                    
      // Reading Progress (for long articles)                                                                                                                                                                       
      EstimatedTime int       `json:"estimatedTime"` // in minutes                                                                                                                                                  
                                                                                                                                                                                                                    
      // Rich Media                                                                                                                                                                                                 
      VideoURL      string    `json:"videoUrl" gorm:"size:500"`                                                                                                                                                     
                                                                                                                                                                                                                    
      // Call to Action                                                                                                                                                                                             
      CTAText       string    `json:"ctaText" gorm:"size:100"`                                                                                                                                                      
      CTALink       string    `json:"ctaLink" gorm:"size:500"`                                                                                                                                                      
                                                                                                                                                                                                                    
      // Analytics                                                                                                                                                                                                  
      AvgReadTime   float64   `json:"avgReadTime"` // actual average                                                                                                                                                
      BounceRate    float64   `json:"bounceRate"`                                                                                                                                                                   
      CompletionRate float64   `json:"completionRate"`  
}                        

 // Category (hierarchical if needed)                                                                                                                                                                              
  type Category struct {                                                                                                                                                                                            
      ID          uuid.UUID    `json:"id"         gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`                                                                                                            
      Name        string       `json:"name"       gorm:"size:100;not null"`                                                                                                                                         
      Slug        string       `json:"slug"       gorm:"size:100;uniqueIndex;not null"`                                                                                                                             
      Description string       `json:"description" gorm:"size:255"`                                                                                                                                                 
      ParentID    *uuid.UUID   `json:"parentId"   gorm:"type:uuid"`                                                                                                                                                 
      Parent      *Category    `json:"parent,omitempty" gorm:"foreignKey:ParentID"`                                                                                                                                 
      Order       int          `json:"order"      gorm:"default:0"`                                                                                                                                                 
      CreatedAt   time.Time    `json:"createdAt"`                                                                                                                                                                   
      UpdatedAt   time.Time    `json:"updatedAt"`                                                                                                                                                                   
  }                                                                                                                                                                                                                 
                                                                                                                                                                                                                    
  // ArticleTag (many-to-many junction if not using PostgreSQL arrays)                                                                                                                                              
  type ArticleTag struct {                                                                                                                                                                                          
      ID         uuid.UUID  `json:"id"      gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`                                                                                                                  
      ArticleID  uuid.UUID  `json:"articleId" gorm:"type:uuid;uniqueIndex:idx_article_tag"`                                                                                                                         
      TagID      uuid.UUID  `json:"tagId"   gorm:"type:uuid;uniqueIndex:idx_article_tag"`                                                                                                                           
      Tag        *Tag       `json:"tag,omitempty" gorm:"foreignKey:TagID"`                                                                                                                                          
  }                                                                                                                                                                                                                 
                                                                                                                                                                                                                    
  type Tag struct {                                                                                                                                                                                                 
      ID          uuid.UUID   `json:"id"         gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`                                                                                                             
      Name        string      `json:"name"       gorm:"size:50;uniqueIndex;not null"`                                                                                                                               
      Slug        string      `json:"slug"       gorm:"size:50;uniqueIndex;not null"`                                                                                                                               
      Description string      `json:"description" gorm:"size:255"`                                                                                                                                                  
      UsageCount  int         `json:"usageCount" gorm:"default:0"`                                                                                                                                                  
      CreatedAt   time.Time   `json:"createdAt"`                                                                                                                                                                    
  }                                                                                                                                                                                                                 
                                                                                                                                                                                                                    
  // RelatedArticle (explicit relationships)                                                                                                                                                                        
  type RelatedArticle struct {                                                                                                                                                                                      
      ID              uuid.UUID  `json:"id"        gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`                                                                                                           
      ArticleID       uuid.UUID  `json:"articleId" gorm:"type:uuid;uniqueIndex:idx_related"`                                                                                                                        
      RelatedArticleID uuid.UUID `json:"relatedId" gorm:"type:uuid;uniqueIndex:idx_related"`                                                                                                                        
      RelationshipType string    `json:"type"      gorm:"size:20"`  // "related", "previous", "next", "sponsored"                                                                                                   
  }