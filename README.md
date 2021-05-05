# Sertifika Doğrula

1 saatte bir Google Sheets içindeki sertifikalar memory'e alınıyor.

Gelen request içindeki JWT token parçalanıyor. Eğer doğrulama başarılı ise 
memory içinden kontrol ediliyor ve doğrulanırsa aşağıdaki şekilde yanıt dönüyor.

Örnek Request:

```
POST https://URL.com/api/certificate_verification

{
    "token": "eeeeeeeeeeeeeeee",
    "qr_code": "Ae1epOlMn"
}
```

Başarılı yanıt:
```json
{
  "status": "verified",
  "qr_code": "Ae1epOlMn",
  "certificate_name": "İş Sağlığı ve Ergonomi",
  "first_name": "Ay***",
  "last_name": "Çot***"
}
```

Başarısız yanıt:
```json
{
  "status": "not_verified",
  "qr_code": "Ae1epOlMn",
  "certificate_name": null,
  "first_name": null,
  "last_name": null
}
```

## ENV Variables

- JWT_SECRET
- Sheet ID
- Page ID
