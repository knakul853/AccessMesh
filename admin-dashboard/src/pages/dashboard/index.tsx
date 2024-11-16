import React from 'react';
import {
  Box,
  Card,
  CardContent,
  Grid,
  Typography,
} from '@mui/material';
import DashboardLayout from '@/components/Layout/DashboardLayout';
import { Security, Policy, Person } from '@mui/icons-material';
import useSWR from 'swr';
import api from '@/lib/api';

const fetcher = (url: string) => api.get(url).then((res) => res.data);

interface Stats {
  totalRoles: number;
  totalPolicies: number;
  totalUsers: number;
}

export default function DashboardPage() {
  const { data: stats, error } = useSWR<Stats>('/api/stats', fetcher);

  const statsCards = [
    {
      title: 'Total Roles',
      value: stats?.totalRoles || 0,
      icon: <Security sx={{ fontSize: 40 }} />,
      color: '#1976d2',
    },
    {
      title: 'Total Policies',
      value: stats?.totalPolicies || 0,
      icon: <Policy sx={{ fontSize: 40 }} />,
      color: '#2e7d32',
    },
    {
      title: 'Total Users',
      value: stats?.totalUsers || 0,
      icon: <Person sx={{ fontSize: 40 }} />,
      color: '#ed6c02',
    },
  ];

  if (error) {
    return (
      <DashboardLayout>
        <Typography color="error">Failed to load dashboard stats</Typography>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" sx={{ mb: 4 }}>
          Dashboard
        </Typography>

        <Grid container spacing={3}>
          {statsCards.map((stat) => (
            <Grid item xs={12} sm={4} key={stat.title}>
              <Card>
                <CardContent>
                  <Box
                    sx={{
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'space-between',
                    }}
                  >
                    <Box>
                      <Typography
                        variant="h6"
                        sx={{ color: 'text.secondary', mb: 1 }}
                      >
                        {stat.title}
                      </Typography>
                      <Typography variant="h4">{stat.value}</Typography>
                    </Box>
                    <Box
                      sx={{
                        backgroundColor: `${stat.color}15`,
                        borderRadius: '50%',
                        p: 2,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center',
                      }}
                    >
                      {React.cloneElement(stat.icon, { sx: { color: stat.color } })}
                    </Box>
                  </Box>
                </CardContent>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Box>
    </DashboardLayout>
  );
}
