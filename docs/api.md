# Authentication and user management

All protected requests use `Authorization: Bearer <token>`. Login tokens are valid for 30 minutes.

| Method | Path | Access | Description |
|---|---|---|---|
| POST | `/auth/login` | Public | Sign in with email and password; returns `{"token":"...","user":{...}}` and sets the HttpOnly auth cookie |
| GET | `/auth/me` | Authenticated | Return the signed-in user's identity and role |
| PATCH | `/auth/profile` | Authenticated | Update the signed-in user's name, email, employee number, and phone |
| PATCH | `/auth/password` | Authenticated | Change password using current password and matching new-password confirmation |
| GET | `/users` | Administrator | List users |
| POST | `/users` | Administrator | Add a user |
| PATCH | `/users/{id}/role` | Administrator | Change a user's role |

Create a user with `full_name`, `email`, `password`, `confirm_password`, and `role` (`administrator` or `clerk`). `employee_number` and `phone_number` are optional. The initial password is stored as a bcrypt hash. An administrator cannot demote their own account, and the final administrator cannot be demoted.

## Documents

| Method | Path | Access | Description |
|---|---|---|---|
| GET | `/pdf` | Public | List uploaded documents |
| GET | `/pdf/{id}/file` | Public | View an uploaded PDF inline |
| POST | `/pdf/submit` | Public | Upload a PDF |
| DELETE | `/pdf/{id}` | Authenticated | Delete a document and its stored PDF file |