# Authentication and user management

All protected requests use `Authorization: Bearer <token>`. Login tokens are valid for 30 minutes.

| Method | Path | Access | Description |
|---|---|---|---|
| POST | `/auth/login` | Public | Sign in with email and password; returns a token and user |
| GET | `/auth/me` | Authenticated | Return the signed-in user's identity and role |
| GET | `/users` | Authenticated | List users |
| POST | `/users` | Administrator | Add a user |
| PATCH | `/users/{id}/role` | Administrator | Change a user's role |

Create a user with `full_name`, `email`, `password`, and `role` (`administrator` or `clerk`). The initial password is stored as a bcrypt hash. An administrator cannot demote their own account, and the final administrator cannot be demoted.