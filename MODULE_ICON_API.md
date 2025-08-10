# Module Icon Upload API Documentation

## **Upload Module Icon**

### **Endpoint**
```
POST /module/{id}/upload-icon
```

### **Description**
Upload an icon file for a specific module. Only admin users can perform this action.

### **Headers**
- `Authorization: Bearer {token}` (Admin token required)
- `Content-Type: multipart/form-data`

### **Parameters**
- `id` (path parameter): The ID of the module to upload the icon for

### **Body**
- `icon` (file): The icon file to upload
  - **Supported formats**: `.jpg`, `.jpeg`, `.png`, `.svg`, `.gif`
  - **Maximum size**: 5MB
  - **Required**: Yes

### **Response**

#### **Success Response (200 OK)**
```json
{
  "message": "Icon uploaded successfully",
  "data": {
    "ID": 1,
    "title": "Introduction to Programming",
    "description": "Learn the basics of programming",
    "offset_x": 100.5,
    "offset_y": 200.0,
    "icon": "uploads/icons/programming_icon_1722518388.png",
    "thumbnail": "",
    "created_at": "2025-08-01T13:06:28Z",
    "updated_at": "2025-08-01T13:06:28Z"
  },
  "icon_url": "http://localhost:3000/uploads/icons/programming_icon_1722518388.png"
}
```

#### **Error Responses**

**400 Bad Request** - Invalid module ID
```json
{
  "error": "Invalid module ID"
}
```

**404 Not Found** - Module not found
```json
{
  "error": "Module not found",
  "message": "Module with the given ID does not exist"
}
```

**400 Bad Request** - File upload error
```json
{
  "error": "File type .txt not allowed. Allowed types: [.jpg .jpeg .png .svg .gif]",
  "message": "Failed to upload icon"
}
```

**413 Request Entity Too Large** - File too large
```json
{
  "error": "file size exceeds limit of 5242880 bytes",
  "message": "Failed to upload icon"
}
```

**401 Unauthorized** - Missing or invalid token
```json
{
  "error": "Unauthorized",
  "message": "Authorization header not provided"
}
```

**403 Forbidden** - Not admin user
```json
{
  "error": "Access denied",
  "message": "Admin access required"
}
```

---

## **Delete Module Icon**

### **Endpoint**
```
DELETE /module/{id}/delete-icon
```

### **Description**
Delete the icon file for a specific module. Only admin users can perform this action.

### **Headers**
- `Authorization: Bearer {token}` (Admin token required)

### **Parameters**
- `id` (path parameter): The ID of the module to delete the icon from

### **Response**

#### **Success Response (200 OK)**
```json
{
  "message": "Icon deleted successfully",
  "data": {
    "ID": 1,
    "title": "Introduction to Programming",
    "description": "Learn the basics of programming",
    "offset_x": 100.5,
    "offset_y": 200.0,
    "icon": "",
    "thumbnail": "",
    "created_at": "2025-08-01T13:06:28Z",
    "updated_at": "2025-08-01T13:10:15Z"
  }
}
```

#### **Error Responses**

**400 Bad Request** - Invalid module ID
```json
{
  "error": "Invalid module ID"
}
```

**404 Not Found** - Module not found
```json
{
  "error": "Module not found",
  "message": "Module with the given ID does not exist"
}
```

**400 Bad Request** - No icon to delete
```json
{
  "error": "No icon to delete",
  "message": "Module does not have an icon"
}
```

---

## **Access Uploaded Files**

### **Endpoint**
```
GET /uploads/icons/{filename}
```

### **Description**
Access uploaded icon files directly. This endpoint serves static files.

### **Example**
```
GET http://localhost:3000/uploads/icons/programming_icon_1722518388.png
```

---

## **Example Usage with cURL**

### **Upload Icon**
```bash
curl -X POST \
  http://localhost:3000/module/1/upload-icon \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -F "icon=@/path/to/your/icon.png"
```

### **Delete Icon**
```bash
curl -X DELETE \
  http://localhost:3000/module/1/delete-icon \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### **View Icon**
```bash
curl http://localhost:3000/uploads/icons/programming_icon_1722518388.png
```

---

## **Example Usage with JavaScript/Fetch**

### **Upload Icon**
```javascript
const formData = new FormData();
formData.append('icon', fileInput.files[0]);

fetch('http://localhost:3000/module/1/upload-icon', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer ' + adminToken
  },
  body: formData
})
.then(response => response.json())
.then(data => {
  console.log('Icon uploaded:', data);
  console.log('Icon URL:', data.icon_url);
})
.catch(error => console.error('Error:', error));
```

### **Delete Icon**
```javascript
fetch('http://localhost:3000/module/1/delete-icon', {
  method: 'DELETE',
  headers: {
    'Authorization': 'Bearer ' + adminToken
  }
})
.then(response => response.json())
.then(data => {
  console.log('Icon deleted:', data);
})
.catch(error => console.error('Error:', error));
```

---

## **File Storage**

- **Storage Location**: `./uploads/icons/`
- **File Naming**: Original filename with timestamp suffix for uniqueness
- **Example**: `programming_icon_1722518388.png`
- **Access**: Files are accessible via `/uploads/icons/{filename}` endpoint

---

## **Security Notes**

1. **Admin Only**: Both upload and delete operations require admin privileges
2. **File Validation**: File type and size are validated before upload
3. **Unique Filenames**: Timestamp is added to prevent filename conflicts
4. **Old File Cleanup**: Previous icon is automatically deleted when uploading new one
5. **Safe File Serving**: Static files are served safely through Gin's static handler

---

## **Configuration**

### **File Upload Settings**
- **Max File Size**: 5MB (configurable in `utils/file_upload.go`)
- **Allowed Types**: `.jpg`, `.jpeg`, `.png`, `.svg`, `.gif`
- **Upload Directory**: `./uploads/icons/`

### **Customize Settings**
You can modify the upload configuration in `internal/utils/file_upload.go`:

```go
func DefaultImageUploadConfig() FileUploadConfig {
    return FileUploadConfig{
        UploadDir:      "./uploads/icons",
        MaxFileSize:    5 * 1024 * 1024, // 5MB
        AllowedTypes:   []string{".jpg", ".jpeg", ".png", ".svg", ".gif"},
        GenerateUnique: true,
    }
}
```
