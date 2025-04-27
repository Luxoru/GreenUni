import { Redirect } from 'expo-router';
import { useAuth } from './auth/AuthContext';
import { View, Text } from 'react-native';

export default function Index() {
  const { isLoggedIn } = useAuth();

  if (isLoggedIn === null) {
    return (
      <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
        <Text>Loading...</Text>
      </View>
    );
  }

  return isLoggedIn ? <Redirect href="/explore" /> : <Redirect href="/auth/login" />;
}