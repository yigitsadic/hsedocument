# Sertifika Doğrula

1 saatte bir Google Sheets içindeki sertifikalar memory'e alınıyor.

Gelen request içindeki JWT token parçalanıyor. Eğer doğrulama başarılı ise 
memory içinden kontrol ediliyor ve doğrulanırsa aşağıdaki şekilde yanıt dönüyor.

Örnek Request:

```
POST https://URL.com/api/certificate_verification

{
    "token": "0c28a727-dae5-4549-87c0-f074e9a40041",
    "qr_code": "Ae1epOlMn"
}
```

Başarılı yanıt:
```json
{
  "status": "verified",
  "qr_code": "0c28a727-dae5-4549-87c0-f074e9a40041",
  "certificate_name": "İş Sağlığı ve Ergonomi",
  "first_name": "Ay***",
  "last_name": "***oy",
  "certificate_created_at": "2021-05-03"
}
```

Başarısız yanıt:
```json
{
  "status": "not_verified",
  "qr_code": "Ae1epOlMn",
  "certificate_name": "",
  "first_name": "",
  "last_name": "",
  "certificate_created_at": ""
}
```

## ENV Variables

- JWT_SECRET
- Sheet ID
