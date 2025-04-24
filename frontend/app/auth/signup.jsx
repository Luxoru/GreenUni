// app/auth/signup.jsx
import React, { useState } from 'react';
import { View, Text, TouchableOpacity, TextInput, StyleSheet, Alert } from 'react-native';
import { useRouter } from 'expo-router';

export default function SignUp() {
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [selectedRole, setSelectedRole] = useState('Student');
  const router = useRouter();

  const roles = ['Student', 'Recruiter'];

  const handleSignUp = async () => {
    try {
      const response = await fetch('http://192.168.1.58:8080/api/v1/auth/signup', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          username,
          password,
          email,
          role: selectedRole
        }),
      });

      const data = await response.json();
      if (!response.ok) {
        throw new Error(data.message || 'Signup failed, unknown error occured');
      }

      if(!data.success){
        throw new Error(data.message)
      }

      console.log('Signup success:', data);
      router.push('/auth/login');
    } catch (error) {
      console.error('Signup error:', error);
      Alert.alert('Signup Failed', error.message);
    }
  };

  return (
    <View style={styles.container}>
      <Text style={styles.heading}>Create Your Account</Text>

      <TextInput
        style={styles.input}
        placeholder="Email"
        autoCapitalize="none"
        keyboardType="email-address"
        value={email}
        onChangeText={setEmail}
      />

      <TextInput
        style={styles.input}
        placeholder="Username"
        autoCapitalize="none"
        value={username}
        onChangeText={setUsername}
      />

      <TextInput
        style={styles.input}
        placeholder="Password"
        secureTextEntry
        value={password}
        onChangeText={setPassword}
      />

      <Text style={styles.subHeading}>Select Role:</Text>
      {roles.map((role) => (
        <TouchableOpacity
          key={role}
          style={styles.option}
          onPress={() => setSelectedRole(role)}
        >
          <View style={[styles.radio, selectedRole === role && styles.radioSelected]}>
            {selectedRole === role && <View style={styles.radioDot} />}
          </View>
          <Text style={styles.label}>{role}</Text>
        </TouchableOpacity>
      ))}

      <TouchableOpacity style={styles.button} onPress={handleSignUp}>
        <Text style={styles.buttonText}>Sign Up</Text>
      </TouchableOpacity>

      <Text style={styles.orText}>or</Text>

      <TouchableOpacity style={styles.button} onPress={() => router.push('/auth/login')}>
        <Text style={styles.buttonText}>Already have an account? Log in</Text>
      </TouchableOpacity>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    padding: 24,
    backgroundColor: '#fff',
  },
  heading: {
    fontSize: 26,
    fontWeight: 'bold',
    marginBottom: 30,
    textAlign: 'center',
  },
  subHeading: {
    fontSize: 16,
    marginTop: 10,
    marginBottom: 10,
  },
  input: {
    height: 44,
    borderWidth: 1,
    borderColor: '#ccc',
    borderRadius: 8,
    paddingHorizontal: 10,
    marginBottom: 12,
  },
  option: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 10,
  },
  radio: {
    height: 20,
    width: 20,
    borderRadius: 10,
    borderWidth: 2,
    borderColor: '#888',
    marginRight: 10,
    justifyContent: 'center',
    alignItems: 'center',
  },
  radioSelected: {
    borderColor: '#007AFF',
  },
  radioDot: {
    height: 10,
    width: 10,
    borderRadius: 5,
    backgroundColor: '#007AFF',
  },
  label: {
    fontSize: 16,
  },
  button: {
    backgroundColor: '#007AFF',
    borderRadius: 8,
    paddingVertical: 12,
    marginTop: 5,
  },
  buttonText: {
    color: '#fff',
    textAlign: 'center',
    fontSize: 16,
  },
  orText: {
    textAlign: 'center',
    marginVertical: 14,
    fontSize: 16,
  },
  loginText: {
    textAlign: 'center',
    fontSize: 16,
    color: '#007AFF',
  },
});
