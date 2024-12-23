import { useState, useEffect } from 'react';
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
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  IconButton,
  CircularProgress
} from '@mui/material';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import AddIcon from '@mui/icons-material/Add';
import { API_BASE_URL } from '../config/api';
import { useRouter } from 'next/router';
import DashboardLayout from '@/components/Layout/DashboardLayout';

interface Policy {
  _id: string;
  role: string;
  resource: string;
  action: string;
  conditions: {
    ip_range: string[];
    time_range: string[];
  };
}

export default function Policies() {
  const [policies, setPolicies] = useState<Policy[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [openDialog, setOpenDialog] = useState(false);
  const [editingPolicy, setEditingPolicy] = useState<Policy | null>(null);
  const [formData, setFormData] = useState({
    role: '',
    resource: '',
    action: '',
    conditions: {
      ip_range: [''],
      time_range: ['']
    }
  });
  const router = useRouter();

  const fetchPolicies = async () => {
    try {
      const token = localStorage.getItem('accessToken');
      if (!token) {
        router.push(`/login?returnTo=${encodeURIComponent('/policies')}`);
        return;
      }

      const response = await fetch(`${API_BASE_URL}/api/v1/policies`, {
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `Failed to fetch policies: ${response.status}`);
      }

      const data = await response.json();
      setPolicies(data);
      setError(''); // Clear any existing errors
    } catch (err) {
      console.error('Error fetching policies:', err);
      setError(err instanceof Error ? err.message : 'Failed to fetch policies');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPolicies();
  }, []);

  const handleOpenDialog = (policy?: Policy) => {
    if (policy) {
      setEditingPolicy(policy);
      setFormData({
        role: policy.role,
        resource: policy.resource,
        action: policy.action,
        conditions: {
          ip_range: policy.conditions.ip_range.length > 0 ? policy.conditions.ip_range : [''],
          time_range: policy.conditions.time_range.length > 0 ? policy.conditions.time_range : ['']
        }
      });
    } else {
      setEditingPolicy(null);
      setFormData({
        role: '',
        resource: '',
        action: '',
        conditions: {
          ip_range: [''],
          time_range: ['']
        }
      });
    }
    setOpenDialog(true);
  };

  const handleCloseDialog = () => {
    setOpenDialog(false);
    setEditingPolicy(null);
    setFormData({
      role: '',
      resource: '',
      action: '',
      conditions: {
        ip_range: [''],
        time_range: ['']
      }
    });
  };

  const handleSubmit = async () => {
    try {
      const token = localStorage.getItem('accessToken');
      if (!token) {
        router.push(`/login?returnTo=${encodeURIComponent('/policies')}`);
        return;
      }

      // Log the request payload
      console.log('Submitting policy with data:', formData);

      const url = editingPolicy
        ? `${API_BASE_URL}/api/v1/policies/${editingPolicy._id}`
        : `${API_BASE_URL}/api/v1/policies`;

      console.log('Making request to:', url);

      const requestBody = {
        ...formData,
        conditions: {
          ip_range: formData.conditions.ip_range.filter(ip => ip.trim() !== ''),
          time_range: formData.conditions.time_range.filter(time => time.trim() !== '')
        }
      };

      console.log('Request body:', requestBody);

      const response = await fetch(url, {
        method: editingPolicy ? 'PUT' : 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(requestBody)
      });

      console.log('Response status:', response.status);
      
      const responseText = await response.text();
      console.log('Raw response:', responseText);

      let data;
      try {
        data = responseText ? JSON.parse(responseText) : {};
      } catch (e) {
        console.error('Error parsing response:', e);
        throw new Error('Invalid response from server');
      }

      if (!response.ok) {
        throw new Error(data.error || `Failed to ${editingPolicy ? 'update' : 'create'} policy: ${response.status}`);
      }

      console.log('Policy saved successfully:', data);
      
      handleCloseDialog();
      await fetchPolicies(); // Refresh the policies list
      setError(''); // Clear any existing errors
    } catch (err) {
      console.error('Error saving policy:', err);
      setError(err instanceof Error ? err.message : 'Failed to save policy');
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this policy?')) {
      return;
    }

    try {
      const token = localStorage.getItem('accessToken');
      if (!token) {
        router.push(`/login?returnTo=${encodeURIComponent('/policies')}`);
        return;
      }

      const response = await fetch(`${API_BASE_URL}/api/v1/policies/${id}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        }
      });

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}));
        throw new Error(errorData.error || `Failed to delete policy: ${response.status}`);
      }

      await fetchPolicies(); // Refresh the policies list
      setError(''); // Clear any existing errors
    } catch (err) {
      console.error('Error deleting policy:', err);
      setError(err instanceof Error ? err.message : 'Failed to delete policy');
    }
  };

  const handleAddCondition = (type: 'ip_range' | 'time_range') => {
    setFormData(prev => ({
      ...prev,
      conditions: {
        ...prev.conditions,
        [type]: [...prev.conditions[type], '']
      }
    }));
  };

  const handleConditionChange = (type: 'ip_range' | 'time_range', index: number, value: string) => {
    setFormData(prev => ({
      ...prev,
      conditions: {
        ...prev.conditions,
        [type]: prev.conditions[type].map((item, i) => i === index ? value : item)
      }
    }));
  };

  const handleRemoveCondition = (type: 'ip_range' | 'time_range', index: number) => {
    setFormData(prev => ({
      ...prev,
      conditions: {
        ...prev.conditions,
        [type]: prev.conditions[type].filter((_, i) => i !== index)
      }
    }));
  };

  return (
    <DashboardLayout>
      <Box sx={{ mb: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
          <Typography variant="h4">Policies</Typography>
          <Button 
            variant="contained" 
            color="primary"
            onClick={() => handleOpenDialog()}
          >
            Create New Policy
          </Button>
        </Box>

        {loading ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
            <CircularProgress />
          </Box>
        ) : error ? (
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <Typography color="error" gutterBottom>
              {error}
            </Typography>
            <Button variant="contained" onClick={fetchPolicies}>
              Retry
            </Button>
          </Box>
        ) : policies.length === 0 ? (
          <Box sx={{ textAlign: 'center', py: 4 }}>
            <Typography color="textSecondary" gutterBottom>
              No policies found
            </Typography>
            <Button 
              variant="contained" 
              color="primary"
              onClick={() => handleOpenDialog()}
            >
              Create First Policy
            </Button>
          </Box>
        ) : (
          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Role</TableCell>
                  <TableCell>Resource</TableCell>
                  <TableCell>Action</TableCell>
                  <TableCell>Conditions</TableCell>
                  <TableCell>Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {policies.map((policy) => (
                  <TableRow key={policy._id}>
                    <TableCell>{policy.role}</TableCell>
                    <TableCell>{policy.resource}</TableCell>
                    <TableCell>{policy.action}</TableCell>
                    <TableCell>
                      <Box>
                        {policy.conditions.ip_range.length > 0 && (
                          <Box mb={1}>
                            <Typography variant="subtitle2">IP Ranges:</Typography>
                            <ul style={{ margin: 0, paddingLeft: 20 }}>
                              {policy.conditions.ip_range.map((ip, index) => (
                                <li key={index}>{ip}</li>
                              ))}
                            </ul>
                          </Box>
                        )}
                        {policy.conditions.time_range.length > 0 && (
                          <Box>
                            <Typography variant="subtitle2">Time Ranges:</Typography>
                            <ul style={{ margin: 0, paddingLeft: 20 }}>
                              {policy.conditions.time_range.map((time, index) => (
                                <li key={index}>{time}</li>
                              ))}
                            </ul>
                          </Box>
                        )}
                      </Box>
                    </TableCell>
                    <TableCell>
                      <IconButton onClick={() => handleOpenDialog(policy)} color="primary">
                        <EditIcon />
                      </IconButton>
                      <IconButton onClick={() => handleDelete(policy._id)} color="error">
                        <DeleteIcon />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        )}
      </Box>

      <Dialog open={openDialog} onClose={handleCloseDialog} maxWidth="md" fullWidth>
        <DialogTitle>
          {editingPolicy ? 'Edit Policy' : 'Add New Policy'}
        </DialogTitle>
        <DialogContent>
          <Box display="flex" flexDirection="column" gap={2} mt={2}>
            <TextField
              label="Role"
              value={formData.role}
              onChange={(e) => setFormData(prev => ({ ...prev, role: e.target.value }))}
              fullWidth
              required
              helperText="e.g., admin, user, editor"
            />
            <TextField
              label="Resource"
              value={formData.resource}
              onChange={(e) => setFormData(prev => ({ ...prev, resource: e.target.value }))}
              fullWidth
              required
              helperText="e.g., /api/users, /api/posts"
            />
            <TextField
              label="Action"
              value={formData.action}
              onChange={(e) => setFormData(prev => ({ ...prev, action: e.target.value }))}
              fullWidth
              required
              helperText="e.g., read, write, delete"
            />
            
            <Typography variant="h6" mt={2}>IP Range Conditions</Typography>
            {formData.conditions.ip_range.map((ip, index) => (
              <Box key={index} display="flex" gap={1}>
                <TextField
                  label={`IP Range ${index + 1}`}
                  value={ip}
                  onChange={(e) => handleConditionChange('ip_range', index, e.target.value)}
                  fullWidth
                  helperText="e.g., 192.168.1.0/24"
                />
                <IconButton
                  color="error"
                  onClick={() => handleRemoveCondition('ip_range', index)}
                  disabled={formData.conditions.ip_range.length === 1}
                >
                  <DeleteIcon />
                </IconButton>
              </Box>
            ))}
            <Button
              startIcon={<AddIcon />}
              onClick={() => handleAddCondition('ip_range')}
              variant="outlined"
              size="small"
            >
              Add IP Range
            </Button>

            <Typography variant="h6" mt={2}>Time Range Conditions</Typography>
            {formData.conditions.time_range.map((time, index) => (
              <Box key={index} display="flex" gap={1}>
                <TextField
                  label={`Time Range ${index + 1}`}
                  value={time}
                  onChange={(e) => handleConditionChange('time_range', index, e.target.value)}
                  fullWidth
                  helperText="e.g., Mon-Fri 9:00-17:00"
                />
                <IconButton
                  color="error"
                  onClick={() => handleRemoveCondition('time_range', index)}
                  disabled={formData.conditions.time_range.length === 1}
                >
                  <DeleteIcon />
                </IconButton>
              </Box>
            ))}
            <Button
              startIcon={<AddIcon />}
              onClick={() => handleAddCondition('time_range')}
              variant="outlined"
              size="small"
            >
              Add Time Range
            </Button>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDialog}>Cancel</Button>
          <Button onClick={handleSubmit} variant="contained" color="primary">
            {editingPolicy ? 'Update' : 'Create'}
          </Button>
        </DialogActions>
      </Dialog>
    </DashboardLayout>
  );
}
