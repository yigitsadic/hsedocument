# Sertifika Doğrula

1 saatte bir Google Sheets içindeki sertifikalar memory'e alınıyor.

Gelen request içindeki JWT token parçalanıyor. Eğer doğrulama başarılı ise 
memory içinden kontrol ediliyor ve doğrulanırsa aşağıdaki şekilde yanıt dönüyor.

Başarılı yanıt:
```json
{
  "status": "verified",
  "qr_code": "AB12C1ia1w",
  "first_name": "Ay***",
  "last_name": "Çot***"
}
```

Başarısız yanıt:
```json
{
  "status": "not_verified",
  "qr_code": "AB12C1ia1w",
  "first_name": "",
  "last_name": ""
}
```

## ENV Variables

- JWT_SECRET
- Sheet ID
- Page ID
