export interface Role {
  id: string;
  name: string;
  description: string;
  permissions: string[];
}

export interface Policy {
  id: string;
  name: string;
  description: string;
  rules: string[];
  createdAt: string;
  updatedAt: string;
}
