import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Box,
  Chip,
  Typography,
} from '@mui/material';
import { Role } from '@/types';

interface RoleFormProps {
  open: boolean;
  onClose: () => void;
  onSubmit: (role: Partial<Role>) => void;
  role?: Role;
}

export default function RoleForm({ open, onClose, onSubmit, role }: RoleFormProps) {
  const [name, setName] = React.useState(role?.name || '');
  const [description, setDescription] = React.useState(role?.description || '');
  const [permission, setPermission] = React.useState('');
  const [permissions, setPermissions] = React.useState<string[]>(role?.permissions || []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit({
      name,
      description,
      permissions,
    });
  };

  const handleAddPermission = () => {
    if (permission && !permissions.includes(permission)) {
      setPermissions([...permissions, permission]);
      setPermission('');
    }
  };

  const handleRemovePermission = (permToRemove: string) => {
    setPermissions(permissions.filter(p => p !== permToRemove));
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{role ? 'Edit Role' : 'Create New Role'}</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mt: 2 }}>
            <TextField
              label="Name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
              fullWidth
            />
            <TextField
              label="Description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              multiline
              rows={3}
              fullWidth
            />
            <Box>
              <Typography variant="subtitle2" gutterBottom>
                Permissions
              </Typography>
              <Box sx={{ display: 'flex', gap: 1, mb: 1 }}>
                <TextField
                  label="Add Permission"
                  value={permission}
                  onChange={(e) => setPermission(e.target.value)}
                  size="small"
                  fullWidth
                />
                <Button
                  variant="outlined"
                  onClick={handleAddPermission}
                  disabled={!permission}
                >
                  Add
                </Button>
              </Box>
              <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1 }}>
                {permissions.map((perm) => (
                  <Chip
                    key={perm}
                    label={perm}
                    onDelete={() => handleRemovePermission(perm)}
                  />
                ))}
              </Box>
            </Box>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Cancel</Button>
          <Button type="submit" variant="contained" color="primary">
            {role ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
}
