// app/auth/signup.jsx
import React, { useState } from 'react';
import { View, Text, TextInput, Button, StyleSheet } from 'react-native';
import { useRouter } from 'expo-router'; // For navigation

export default function SignUp() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const router = useRouter();

  const handleSignUp = async () => {
    // Logic for signing up (e.g., save to backend or AsyncStorage)
    console.log('Signing up with:', email, password);

    // On success, navigate to the login page or home screen
    router.push('/auth/login');
  };

  const handleGoToLogin = () => {
    // Navigate to login screen
    router.push('/auth/login');
  };

  return (
    <View style={styles.container}>
      <Text style={styles.heading}>Sign Up</Text>
      
      <TextInput
        style={styles.input}
        placeholder="Email"
        value={email}
        onChangeText={setEmail}
      />
      
      <TextInput
        style={styles.input}
        placeholder="Password"
        secureTextEntry
        value={password}
        onChangeText={setPassword}
      />
      
      <Button title="Sign Up" onPress={handleSignUp} />

      <Text style={styles.orText}>or</Text>

      <Button title="Go to Login" onPress={handleGoToLogin} />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    padding: 16,
  },
  heading: {
    fontSize: 24,
    fontWeight: 'bold',
    marginBottom: 20,
    textAlign: 'center',
  },
  input: {
    height: 40,
    borderColor: '#ccc',
    borderWidth: 1,
    marginBottom: 16,
    paddingLeft: 8,
  },
  orText: {
    textAlign: 'center',
    marginVertical: 12,
    fontSize: 16,
  },
});
