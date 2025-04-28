import { View, TouchableOpacity, StyleSheet, Text, BackHandler, Alert } from 'react-native';
import { IconSymbol } from '@/components/ui/IconSymbol';
import { Ionicons } from '@expo/vector-icons';
import { Colors } from '@/constants/Colors';
import { useColorScheme } from '@/hooks/useColorScheme';
import * as SecureStore from 'expo-secure-store';
import { useState, useEffect } from 'react';
import { useRouter } from 'expo-router';

export function CustomTabBar({ state, descriptors, navigation }) {
  const [user, setUser] = useState(null);
  const colorScheme = useColorScheme();
  const tintColor = Colors[colorScheme ?? 'light'].tint;
  const router = useRouter();

  // Fetch user from SecureStore
  useEffect(() => {
    async function getUserFromStorage() {
      try {
        const userStr = await SecureStore.getItemAsync('user');
        
        if (!userStr) {
          console.warn("No user found in SecureStore.");
          return;
        }
        
        const userData = JSON.parse(userStr);
        setUser(userData);
        console.log("USER ROLE FROM SECURESTORE:", userData?.role);
      } catch (error) {
        console.error("Error getting user from SecureStore:", error);
      }
    }
    
    getUserFromStorage();
  }, []);

  // Handle back button for recruiters
  useEffect(() => {
    if (user?.role !== 'Recruiter') return;

    const backAction = () => {
      // If we're on the recruiterExplore page and user is a recruiter
      if (state.routes[state.index].name === 'recruiterExplore') {
        Alert.alert('Exit App', 'Are you sure you want to exit?', [
          {
            text: 'Cancel',
            onPress: () => null,
            style: 'cancel',
          },
          { text: 'Exit', onPress: () => BackHandler.exitApp() },
        ]);
        return true; // Prevent default back behavior
      }
      return false; // Let the default back behavior happen
    };

    const backHandler = BackHandler.addEventListener('hardwareBackPress', backAction);

    return () => backHandler.remove();
  }, [user, state.index, state.routes]);

  // Get visible routes based on user role
  let visibleRoutes = [];
  
  if (user?.role === 'Recruiter') {
    // For Recruiters, show recruiterExplore, matches, and recruiterProfile tabs
    visibleRoutes = state.routes.filter(route => 
      ['recruiterExplore', 'matches', 'recruiterProfile'].includes(route.name)
    );
  } else {
    // For students, show explore, matches, and profile
    visibleRoutes = state.routes.filter(route => 
      ['explore', 'matches', 'profile'].includes(route.name)
    );
  }

  // Ensure correct tab for Recruiters
  useEffect(() => {
    if (user?.role === 'Recruiter' && 
        state.routeNames.includes('recruiterExplore') && 
        !['recruiterExplore', 'matches', 'recruiterProfile'].includes(state.routes[state.index].name)) {
      router.replace('/recruiterExplore');
    }
  }, [user, state.index]);

  return (
    <View style={styles.container}>
      {visibleRoutes.map((route, index) => {
        const { options } = descriptors[route.key];
        const isFocused = state.index === state.routes.indexOf(route);

        const onPress = () => {
          const event = navigation.emit({
            type: 'tabPress',
            target: route.key,
            canPreventDefault: true,
          });

          if (!isFocused && !event.defaultPrevented) {
            // Use replace instead of navigate to prevent adding to history stack
            if (user?.role === 'Recruiter') {
              router.replace(`/${route.name}`);
            } else {
              navigation.navigate(route.name);
            }
          }
        };

        const getIcon = () => {
          switch (route.name) {
            case 'explore':
              return <IconSymbol size={28} name="paperplane.fill" color={isFocused ? tintColor : '#666'} />;
            case 'recruiterExplore':
              return <IconSymbol size={28} name="paperplane.fill" color={isFocused ? tintColor : '#666'} />;
            case 'matches':
              return <Ionicons size={28} name="heart" color={isFocused ? tintColor : '#666'} />;
            case 'profile':
              return <Ionicons size={28} name="person" color={isFocused ? tintColor : '#666'} />;
            case 'recruiterProfile':
              return <Ionicons size={28} name="person" color={isFocused ? tintColor : '#666'} />;
            default:
              return null;
          }
        };

        // Display title
        const getDisplayTitle = () => {
          if (route.name === 'recruiterExplore') {
            return 'Explore';
          }
          if (route.name === 'recruiterProfile') {
            return 'Profile';
          }
          return options.title || route.name;
        };

        return (
          <TouchableOpacity
            key={route.key}
            accessibilityRole="button"
            accessibilityState={isFocused ? { selected: true } : {}}
            accessibilityLabel={options.tabBarAccessibilityLabel}
            testID={options.tabBarTestID}
            onPress={onPress}
            style={styles.tab}
          >
            {getIcon()}
            <Text style={[
              styles.label,
              { color: isFocused ? tintColor : '#666' }
            ]}>
              {getDisplayTitle()}
            </Text>
          </TouchableOpacity>
        );
      })}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flexDirection: 'row',
    backgroundColor: '#fff',
    borderTopWidth: 1,
    borderTopColor: '#eee',
    paddingBottom: 20,
    paddingTop: 10,
  },
  tab: {
    flex: 1,
    alignItems: 'center',
    justifyContent: 'center',
  },
  label: {
    fontSize: 12,
    marginTop: 4,
  },
}); 