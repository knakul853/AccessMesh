import React from 'react';
import {
  Box,
  Button,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
  Alert,
  Snackbar,
} from '@mui/material';
import DashboardLayout from '@/components/Layout/DashboardLayout';
import RoleForm from '@/components/RoleForm';
import { Role } from '@/types';
import useSWR, { mutate } from 'swr';
import api from '@/lib/api';

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export default function RolesPage() {
  const { data: roles, error } = useSWR<Role[]>('/api/v1/roles', fetcher);
  const [openForm, setOpenForm] = React.useState(false);
  const [selectedRole, setSelectedRole] = React.useState<Role | undefined>();
  const [message, setMessage] = React.useState<{ text: string; type: 'success' | 'error' } | null>(null);

  const handleCreateRole = async (roleData: Partial<Role>) => {
    try {
      await api.post('/api/v1/roles', roleData);
      mutate('/api/v1/roles');
      setOpenForm(false);
      setMessage({ text: 'Role created successfully', type: 'success' });
    } catch (error) {
      setMessage({ text: 'Failed to create role', type: 'error' });
    }
  };

  const handleUpdateRole = async (roleData: Partial<Role>) => {
    if (!selectedRole) return;
    try {
      await api.put(`/api/v1/roles/${selectedRole.id}`, roleData);
      mutate('/api/v1/roles');
      setOpenForm(false);
      setSelectedRole(undefined);
      setMessage({ text: 'Role updated successfully', type: 'success' });
    } catch (error) {
      setMessage({ text: 'Failed to update role', type: 'error' });
    }
  };

  const handleDeleteRole = async (roleId: string) => {
    if (!window.confirm('Are you sure you want to delete this role?')) return;
    try {
      await api.delete(`/api/v1/roles/${roleId}`);
      mutate('/api/v1/roles');
      setMessage({ text: 'Role deleted successfully', type: 'success' });
    } catch (error) {
      setMessage({ text: 'Failed to delete role', type: 'error' });
    }
  };

  const handleEdit = (role: Role) => {
    setSelectedRole(role);
    setOpenForm(true);
  };

  const handleCloseForm = () => {
    setOpenForm(false);
    setSelectedRole(undefined);
  };

  if (error) {
    return (
      <DashboardLayout>
        <Typography color="error">Failed to load roles</Typography>
      </DashboardLayout>
    );
  }

  if (!roles) {
    return (
      <DashboardLayout>
        <Typography>Loading...</Typography>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <Box sx={{ mb: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
          <Typography variant="h4">Roles</Typography>
          <Button 
            variant="contained" 
            color="primary"
            onClick={() => setOpenForm(true)}
          >
            Create New Role
          </Button>
        </Box>

        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Description</TableCell>
                <TableCell>Permissions</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {roles.map((role) => (
                <TableRow key={role.id}>
                  <TableCell>{role.name}</TableCell>
                  <TableCell>{role.description}</TableCell>
                  <TableCell>
                    {role.permissions.map((perm, index) => (
                      <React.Fragment key={perm}>
                        {perm}
                        {index < role.permissions.length - 1 ? ', ' : ''}
                      </React.Fragment>
                    ))}
                  </TableCell>
                  <TableCell>
                    <Button
                      size="small"
                      variant="outlined"
                      sx={{ mr: 1 }}
                      onClick={() => handleEdit(role)}
                    >
                      Edit
                    </Button>
                    <Button
                      size="small"
                      variant="outlined"
                      color="error"
                      onClick={() => handleDeleteRole(role.id)}
                    >
                      Delete
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Box>

      <RoleForm
        open={openForm}
        onClose={handleCloseForm}
        onSubmit={selectedRole ? handleUpdateRole : handleCreateRole}
        role={selectedRole}
      />

      <Snackbar
        open={!!message}
        autoHideDuration={6000}
        onClose={() => setMessage(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert 
          onClose={() => setMessage(null)} 
          severity={message?.type || 'info'}
          sx={{ display: message ? 'flex' : 'none' }}
        >
          {message?.text}
        </Alert>
      </Snackbar>
    </DashboardLayout>
  );
}
