export interface User {
  id: string;
  username: string;
  email: string;
  full_name: string;
  role?: string;
  permissions?: string[];
}

export interface LoginResponse {
  success: {
    status: number;
    data: {
      accessToken: string;
      refreshToken: string;
      user: User;
    };
  };
}
