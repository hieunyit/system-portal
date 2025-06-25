export interface User {
  id: string;
  username: string;
  email: string;
  full_name: string;
  role?: string;
  permissions?: string[];
}

export interface LoginResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}
