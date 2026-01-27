# Backend API Contract

**Authentication**
- Header: `Authorization: Bearer <email>|<role>`
- Example: `Authorization: Bearer faculty@citchennai.net|faculty`

**Base URL**: `http://localhost:8080` (assumed, need to verify server port, usually 8080)

## Endpoints

### Faculty
| Endpoint | Method | Body | Description |
|----------|--------|------|-------------|
| `/certificates/upload` | POST | `{ "certificates": [ { "drive_link": "...", "register_number": "...", "section": "...", "student_name": "...", "uploaded_by": "...", "uploaded_at": "..." } ] }` | Upload up to 10 certificates |
| `/certificates/pending-review` | GET | Query: `limit` (default 50) | Get ML-verified certificates waiting for faculty review |
| `/certificates/review` | POST | `{ "certificate_id": "...", "status": "LEGIT" \| "NOT_LEGIT", "is_legit": boolean }` | Submit manual review |

### Shared / Debug
| Endpoint | Method | Body | Description |
|----------|--------|------|-------------|
| `/faculty/certificate/verify` | POST | `{ "certificate_id": "..." }` | Trigger mock ML verification (Async simulation) |

### HOD
| Endpoint | Method | Query Params | Description |
|----------|--------|--------------|-------------|
| `/hod/dashboard` | GET | - | Dashboard overview |
| `/hod/faculty/students` | GET | `faculty_id` | Student stats by faculty |
| `/hod/student/certificates` | GET | `reg_no` | List certificates for a student |
| `/hod/export/certificates/section` | GET | `section` | Export excel |
| `/hod/export/certificates/student` | GET | `reg_no` | Export excel |

## Models & Enums

**MLStatus**: `PENDING`, `VERIFIED`, `DUPLICATE`
**FacultyStatus**: `PENDING`, `LEGIT`, `NOT_LEGIT`

**Certificate Object**:
```json
{
  "ID": "uuid",
  "DriveLink": "url",
  "RegisterNumber": "string",
  "Section": "string",
  "StudentName": "string",
  "UploadedBy": "string",
  "UploadedAt": "timestamp",
  "MLStatus": "MLStatus",
  "FacultyStatus": "FacultyStatus",
  "IsLegit": boolean,
  "MLScore": number,
  "Archived": boolean
}
```
