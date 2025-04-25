import React, { useState } from 'react';
import { Modal,View, Text, TextInput, Button, StyleSheet } from 'react-native';
import { useRouter } from 'expo-router';
import { useAuth } from '@/app/auth/AuthContext';
import { Alert } from 'react-native';
import * as SecureStore from 'expo-secure-store';



export default function Login() {
  const { login } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const router = useRouter();
  const [showErrorModal, setShowErrorModal] = useState(false);
  const [errorMessage, setErrorMessage] = useState('');


  const handleLogin = async () => {
    try {
      const response = await fetch('http://192.168.1.58:8080/api/v1/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          username: email,
          password: password,
        }),
      });
  
      if (!response.ok) {
        const errorText = await response.message();
        throw new Error(`HTTP error! Status: ${response.status}, Message: ${errorText}`);
      }
  
      const data = await response.json();
  
      if (!data.success) {
        Alert.alert('Login Failed', data.message.trim());
        return;
      }
      console.log(data.info.user)
      await login(data.info.token)
      await SecureStore.setItemAsync('user', JSON.stringify(data.info.user));
  
      router.push('/explore');
    } catch (error) {
      console.error('Login error:', error);
      Alert.alert('Login Failed', 'Unable to login. Please check your credentials or try again later.');
    }
  };
  
  
  return (
    <View style={styles.container}>
      <Text style={styles.heading}>Login</Text>

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

      <Button title="Login" onPress={handleLogin} />
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
});