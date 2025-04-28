import { Tabs, useRouter } from 'expo-router';
import React, { useEffect, useState } from 'react';
import { Platform } from 'react-native';

import { IconSymbol } from '@/components/ui/IconSymbol';
import { Ionicons } from '@expo/vector-icons';
import { Colors } from '@/constants/Colors';
import { useColorScheme } from '@/hooks/useColorScheme';
import { CustomTabBar } from '@/components/CustomTabBar';
import * as SecureStore from 'expo-secure-store';

export default function TabLayout() {
  const colorScheme = useColorScheme();
  const router = useRouter();
  const [isLoggedIn, setIsLoggedIn] = useState(null);
  const [user, setUser] = useState(null);

  // Check if user is logged in
  useEffect(() => {
    async function checkAuth() {
      try {
        const userStr = await SecureStore.getItemAsync('user');
        
        if (!userStr) {
          console.log("No user found - redirecting to login");
          setIsLoggedIn(false);
          return;
        }
        
        const userData = JSON.parse(userStr);
        setUser(userData);
        setIsLoggedIn(true);
        console.log("User authenticated:", userData?.role);
      } catch (error) {
        console.error("Error checking auth:", error);
        setIsLoggedIn(false);
      }
    }
    
    checkAuth();
  }, []);

  // Redirect to login if not logged in
  useEffect(() => {
    if (isLoggedIn === false) {
      router.replace('/auth/signup');
    }
  }, [isLoggedIn]);

  // Optional loading state while checking auth
  if (isLoggedIn === null) {
    return null; // or splash screen
  }

  return (
    <Tabs
      initialRouteName="explore"
      tabBar={(props) => <CustomTabBar {...props} />}
      screenOptions={{
        headerShown: false,
      }}>
      <Tabs.Screen
        name="explore"
        options={{
          title: 'Explore',
        }}
      />
      <Tabs.Screen
        name="matches"
        options={{
          title: 'Matches',
        }}
      />
      <Tabs.Screen
        name="recruiterExplore"
        options={{
          title: 'Recruiter',
        }}
      />
      <Tabs.Screen
        name="profile"
        options={{
          title: 'Profile',
        }}
      />
      <Tabs.Screen
        name="recruiterProfile"
        options={{
          title: 'Profile',
        }}
      />
    </Tabs>
  );
}