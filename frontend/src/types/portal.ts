export interface PortalUser {
  id: string;
  username: string;
  email: string;
  full_name: string;
}

export interface PortalGroup {
  id: string;
  name: string;
  display_name: string;
}

export interface Permission {
  id: string;
  resource: string;
  action: string;
}
