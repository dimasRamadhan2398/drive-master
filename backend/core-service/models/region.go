package models

type Province struct {                                                                                                                                                                                                   
      ID   uint   `json:"id" gorm:"primaryKey"`                                                                                                                                                                            
      Name string `json:"name" gorm:"size:100;not null"`                                                                                                                                                                   
  }                                                                                                                                                                                                                        
                                                                                                                                                                                                                           
  type Regency struct {                                                                                                                                                                                                    
      ID         uint   `json:"id" gorm:"primaryKey"`                                                                                                                                                                      
      ProvinceID uint   `json:"provinceId" gorm:"not null"`                                                                                                                                                                
      Name       string `json:"name" gorm:"size:100;not null"`
      Type       string `json:"type" gorm:"size:100;not null"`                                                                                                                                          
  }                                                                                                                                                                                                                        
                                                                                                                                                                                                                           
  type District struct {                                                                                                                                                                                                   
      ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`                                                                                                                                                                       
      Name      string `json:"name" gorm:"size:100;not null"`
      RegencyID uint   `json:"regencyId" gorm:"not null"`                                                                                                                                                                                                                                                                                                                   
  } 