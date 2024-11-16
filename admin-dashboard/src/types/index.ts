export interface Role {
  id: string;
  name: string;
  description: string;
  permissions: string[];
  createdAt: string;
  updatedAt: string;
}

export interface Policy {
  id: string;
  name: string;
  description: string;
  effect: 'allow' | 'deny';
  actions: string[];
  resources: string[];
  conditions?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

export interface User {
  id: string;
  username: string;
  email: string;
  roles: Role[];
  createdAt: string;
  updatedAt: string;
}

export interface ApiResponse<T> {
  data: T;
  message: string;
  success: boolean;
}
