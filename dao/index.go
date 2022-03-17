package dao

// import (
//     "reflect"
//     "maintainman/database"

//     "gorm.io/gorm"
// )

// func GetAll(searchKeys map[string]interface{}, orderBy string, offset, limit int) *gorm.DB {
//     DB := database.DB

//     if len(searchKeys) > 0 {
//         for k, v := range searchKeys {
//             if reflect.TypeOf(v).Name() == "string" && v != "" {
//                 DB.Where(k+" = ?", v)
//             } else {
//                 DB.Where(k+" = ?", v)
//             }
//         }
//     }

//     if len(orderBy) > 0 {
//         DB.Order(orderBy + " desc")
//     } else {
//         DB.Order("created_at desc")
//     }

//     if offset > 0 {
//         DB.Offset(offset - 1)
//     }

//     if limit > 0 {
//         DB.Limit(limit)
//     }

//     return DB
// }
