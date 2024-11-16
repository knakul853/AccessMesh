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
} from '@mui/material';
import DashboardLayout from '@/components/Layout/DashboardLayout';
import { Role } from '@/types';
import useSWR from 'swr';
import api from '@/lib/api';

const fetcher = (url: string) => api.get(url).then((res) => res.data);

export default function RolesPage() {
  const { data: roles, error } = useSWR<Role[]>('/api/roles', fetcher);

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
          <Button variant="contained" color="primary">
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
                <TableCell>Created At</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {roles.map((role) => (
                <TableRow key={role.id}>
                  <TableCell>{role.name}</TableCell>
                  <TableCell>{role.description}</TableCell>
                  <TableCell>{role.permissions.join(', ')}</TableCell>
                  <TableCell>
                    {new Date(role.createdAt).toLocaleDateString()}
                  </TableCell>
                  <TableCell>
                    <Button
                      size="small"
                      variant="outlined"
                      sx={{ mr: 1 }}
                      onClick={() => {
                        // Handle edit
                      }}
                    >
                      Edit
                    </Button>
                    <Button
                      size="small"
                      variant="outlined"
                      color="error"
                      onClick={() => {
                        // Handle delete
                      }}
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
    </DashboardLayout>
  );
}
