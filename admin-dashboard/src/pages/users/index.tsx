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
  CircularProgress,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
} from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import DashboardLayout from '@/components/Layout/DashboardLayout';
import api from '@/lib/api';

interface User {
  id: string;
  username: string;
  email: string;
  role: string;
}

const UsersPage = () => {
  const [users, setUsers] = React.useState<User[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);
  const [message, setMessage] = React.useState<{ text: string; type: 'success' | 'error' } | null>(null);
  const [editDialogOpen, setEditDialogOpen] = React.useState(false);
  const [selectedUser, setSelectedUser] = React.useState<User | null>(null);
  const [editedUser, setEditedUser] = React.useState<Partial<User>>({});

  const fetchUsers = React.useCallback(async () => {
    try {
      console.log('Fetching users...');
      const response = await api.get('/api/v1/users');
      console.log('Users response:', response.data);
      const userData = Array.isArray(response.data) ? response.data : response.data.users || [];
      setUsers(userData);
      setError(null);
    } catch (err) {
      console.error('Error fetching users:', err);
      setError('Failed to fetch users');
      setMessage({ text: 'Failed to fetch users', type: 'error' });
    } finally {
      setLoading(false);
    }
  }, []);

  React.useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const handleEditClick = (user: User) => {
    setSelectedUser(user);
    setEditedUser(user);
    setEditDialogOpen(true);
  };

  const handleEditClose = () => {
    setEditDialogOpen(false);
    setSelectedUser(null);
    setEditedUser({});
  };

  const handleEditSave = async () => {
    try {
      await api.put(
        `/api/v1/users/${selectedUser?.id}`,
        editedUser,
      );
      await fetchUsers();
      handleEditClose();
      setMessage({ text: 'User updated successfully', type: 'success' });
    } catch (err) {
      console.error('Error updating user:', err);
      setError('Failed to update user');
      setMessage({ text: 'Failed to update user', type: 'error' });
    }
  };

  const handleDeleteClick = async (userId: string) => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      try {
        await api.delete(`/api/v1/users/${userId}`);
        await fetchUsers();
        setMessage({ text: 'User deleted successfully', type: 'success' });
      } catch (err) {
        console.error('Error deleting user:', err);
        setError('Failed to delete user');
        setMessage({ text: 'Failed to delete user', type: 'error' });
      }
    }
  };

  const renderContent = () => {
    if (error) {
      return (
        <Box sx={{ textAlign: 'center', py: 4 }}>
          <Typography color="error" gutterBottom>
            {error}
          </Typography>
          <Button
            variant="contained"
            color="primary"
            onClick={() => {
              setError(null);
              fetchUsers();
            }}
          >
            Retry
          </Button>
        </Box>
      );
    }

    if (loading) {
      return (
        <Box sx={{ textAlign: 'center', py: 4 }}>
          <CircularProgress />
        </Box>
      );
    }

    if (users.length === 0) {
      return (
        <Box sx={{ textAlign: 'center', py: 4 }}>
          <Typography variant="body1" gutterBottom>
            No users found
          </Typography>
        </Box>
      );
    }

    return (
      <>
        <TableContainer component={Paper}>
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Username</TableCell>
                <TableCell>Email</TableCell>
                <TableCell>Role</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {users.map((user) => (
                <TableRow key={user.id}>
                  <TableCell>{user.username}</TableCell>
                  <TableCell>{user.email}</TableCell>
                  <TableCell>{user.role}</TableCell>
                  <TableCell>
                    <IconButton onClick={() => handleEditClick(user)}>
                      <EditIcon />
                    </IconButton>
                    <IconButton onClick={() => handleDeleteClick(user.id)}>
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>

        <Dialog open={editDialogOpen} onClose={handleEditClose}>
          <DialogTitle>Edit User</DialogTitle>
          <DialogContent>
            <Box sx={{ pt: 2 }}>
              <TextField
                fullWidth
                label="Username"
                value={editedUser.username || ''}
                onChange={(e) => setEditedUser({ ...editedUser, username: e.target.value })}
                margin="normal"
              />
              <TextField
                fullWidth
                label="Email"
                value={editedUser.email || ''}
                onChange={(e) => setEditedUser({ ...editedUser, email: e.target.value })}
                margin="normal"
              />
              <TextField
                fullWidth
                label="Role"
                value={editedUser.role || ''}
                onChange={(e) => setEditedUser({ ...editedUser, role: e.target.value })}
                margin="normal"
              />
            </Box>
          </DialogContent>
          <DialogActions>
            <Button onClick={handleEditClose}>Cancel</Button>
            <Button onClick={handleEditSave} variant="contained" color="primary">
              Save
            </Button>
          </DialogActions>
        </Dialog>
      </>
    );
  };

  return (
    <DashboardLayout>
      <Box sx={{ p: 3 }}>
        <Typography variant="h4" gutterBottom>
          Users
        </Typography>
        {renderContent()}
        {message && (
          <Snackbar
            open={!!message}
            autoHideDuration={6000}
            onClose={() => setMessage(null)}
          >
            <Alert severity={message.type} onClose={() => setMessage(null)}>
              {message.text}
            </Alert>
          </Snackbar>
        )}
      </Box>
    </DashboardLayout>
  );
};

export default UsersPage;
