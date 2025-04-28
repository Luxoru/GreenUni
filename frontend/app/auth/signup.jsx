// app/auth/signup.jsx
import React, { useState } from 'react';
import { View, Text, TouchableOpacity, TextInput, StyleSheet, Alert, ScrollView } from 'react-native';
import { useRouter } from 'expo-router';
import { Ionicons } from '@expo/vector-icons';
import config from '../../utils/config';
export default function SignUp() {
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [selectedRole, setSelectedRole] = useState('Student');
  const router = useRouter();


  const [opportunityTitle, setOpportunityTitle] = useState('');
  const [opportunityDescription, setOpportunityDescription] = useState('');
  const [opportunityLocation, setOpportunityLocation] = useState('');
  const [opportunityType, setOpportunityType] = useState('volunteer');
  
  // Media handling forsx multiple URLs and types
  const [mediaItems, setMediaItems] = useState([{ url: '', type: 'Image' }]);
  
  const [opportunityTags, setOpportunityTags] = useState('');

  const roles = ['Student', 'Recruiter'];
  const opportunityTypes = ['volunteer', 'internship', 'job', 'event'];
  const mediaTypes = ['Image', 'Video'];


  const addMediaItem = () => {
    setMediaItems([...mediaItems, { url: '', type: 'Image' }]);
  };

  const removeMediaItem = (index) => {
    if (mediaItems.length > 1) {
      const updatedItems = [...mediaItems];
      updatedItems.splice(index, 1);
      setMediaItems(updatedItems);
    }
  };

  
  const updateMediaItem = (index, field, value) => {
    const updatedItems = [...mediaItems];
    updatedItems[index] = { ...updatedItems[index], [field]: value };
    setMediaItems(updatedItems);
  };

  const handleSignUp = async () => {
    try {
     
      if (selectedRole === 'Recruiter') {
        if (!opportunityTitle || !opportunityDescription || !opportunityLocation) {
          throw new Error('Please fill in all required opportunity fields');
        }
        
        
        const invalidMedia = mediaItems.find(item => !item.url);
        if (invalidMedia) {
          throw new Error('Please fill in all media URLs or remove empty fields');
        }
      }

      
      const signupResponse = await fetch(`${config.apiURL}/api/v1/auth/signup`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          username,
          password,
          email,
          role: selectedRole
        }),
      });

      const signupData = await signupResponse.json();
      if (!signupResponse.ok) {
        throw new Error(signupData.message || 'Signup failed, unknown error occurred');
      }

      if(!signupData.success){
        throw new Error(signupData.message);
      }

      console.log('Signup success:', signupData);

      
      if (selectedRole === 'Recruiter') {
        
        const loginResponse = await fetch(`${config.apiURL}/api/v1/auth/login`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            username,
            password
          }),
        });

        const loginData = await loginResponse.json();
        if (!loginResponse.ok || !loginData.success) {
          throw new Error('Login failed after signup, but account was created');
        }

        console.log('Login success:', loginData);
        
       
        const userUUID = loginData.info.user.uuid;
        const token = loginData.info.token;
        
        if (!userUUID) {
          throw new Error('Could not get user UUID, but account was created');
        }

        
        const tagsArray = opportunityTags.split(',').map(tag => tag.trim()).filter(tag => tag);
        const mediaURLArray = mediaItems.map(item => item.url);
        const mediaTypeArray = mediaItems.map(item => item.type);

        const opportunityResponse = await fetch(`${config.apiURL}/api/v1/opportunities`, {
          method: 'POST',
          headers: { 
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify({
            title: opportunityTitle,
            description: opportunityDescription,
            location: opportunityLocation,
            type: opportunityType,
            author: userUUID,
            mediaType: mediaTypeArray,
            mediaURL: mediaURLArray,
            tags: tagsArray,
            points: Math.floor(Math.random() * 500) + 100 // Random points between 100-600 temp maybe change gamnbling??
          }),
        });

        const opportunityData = await opportunityResponse.json();
        if (!opportunityResponse.ok) {
          console.error('Opportunity creation failed, but user account was created');
        } else {
          console.log('Opportunity created:', opportunityData);
        }
      }

      router.replace('/auth/login');
    } catch (error) {
      console.error('Signup error:', error);
      Alert.alert('Signup Failed', error.message);
    }
  };

  return (
    <ScrollView contentContainerStyle={styles.scrollContainer}>
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
          autoCapitalize="none"
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

        {/* Opportunity form for recruiters */}
        {selectedRole === 'Recruiter' && (
          <View style={styles.opportunityForm}>
            <Text style={styles.formHeading}>Create Your First Opportunity</Text>
            
            <TextInput
              style={styles.input}
              placeholder="Opportunity Title *"
              value={opportunityTitle}
              onChangeText={setOpportunityTitle}
            />
            
            <TextInput
              style={[styles.input, styles.textArea]}
              placeholder="Description *"
              multiline
              numberOfLines={4}
              value={opportunityDescription}
              onChangeText={setOpportunityDescription}
            />
            
            <TextInput
              style={styles.input}
              placeholder="Location *"
              value={opportunityLocation}
              onChangeText={setOpportunityLocation}
            />

            <Text style={styles.subHeading}>Opportunity Type:</Text>
            <View style={styles.typeContainer}>
              {opportunityTypes.map((type) => (
                <TouchableOpacity
                  key={type}
                  style={[
                    styles.typeButton,
                    opportunityType === type && styles.selectedTypeButton
                  ]}
                  onPress={() => setOpportunityType(type)}
                >
                  <Text 
                    style={[
                      styles.typeText,
                      opportunityType === type && styles.selectedTypeText
                    ]}
                  >
                    {type.charAt(0).toUpperCase() + type.slice(1)}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
            
            <Text style={styles.subHeading}>Media *</Text>
            <Text style={styles.mediaSubheading}>Add images or videos for your opportunity</Text>
            
            {mediaItems.map((item, index) => (
              <View key={index} style={styles.mediaItem}>
                <View style={styles.mediaTypeContainer}>
                  {mediaTypes.map(type => (
                    <TouchableOpacity 
                      key={type}
                      style={[
                        styles.mediaTypeButton,
                        item.type === type && styles.selectedMediaTypeButton
                      ]}
                      onPress={() => updateMediaItem(index, 'type', type)}
                    >
                      <Text 
                        style={[
                          styles.mediaTypeText,
                          item.type === type && styles.selectedMediaTypeText
                        ]}
                      >
                        {type}
                      </Text>
                    </TouchableOpacity>
                  ))}
                </View>
                
                <View style={styles.mediaInputContainer}>
                  <TextInput
                    style={styles.mediaInput}
                    placeholder={`${item.type} URL`}
                    value={item.url}
                    onChangeText={(text) => updateMediaItem(index, 'url', text)}
                  />
                  
                  <TouchableOpacity 
                    style={styles.removeButton}
                    onPress={() => removeMediaItem(index)}
                    disabled={mediaItems.length === 1}
                  >
                    <Ionicons 
                      name="close-circle" 
                      size={24} 
                      color={mediaItems.length === 1 ? '#ccc' : '#ff3b30'} 
                    />
                  </TouchableOpacity>
                </View>
              </View>
            ))}
            
            <TouchableOpacity style={styles.addMediaButton} onPress={addMediaItem}>
              <Ionicons name="add-circle" size={20} color="#007AFF" />
              <Text style={styles.addMediaText}>Add Another Media</Text>
            </TouchableOpacity>
            
            <TextInput
              style={styles.input}
              placeholder="Tags (comma separated, e.g. remote, part-time)"
              value={opportunityTags}
              onChangeText={setOpportunityTags}
            />
          </View>
        )}

        <TouchableOpacity style={styles.button} onPress={handleSignUp}>
          <Text style={styles.buttonText}>Sign Up</Text>
        </TouchableOpacity>

        <Text style={styles.orText}>or</Text>

        <TouchableOpacity style={styles.button} onPress={() => router.push('/auth/login')}>
          <Text style={styles.buttonText}>Already have an account? Log in</Text>
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  scrollContainer: {
    flexGrow: 1,
  },
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
    fontWeight: '600',
  },
  mediaSubheading: {
    fontSize: 14,
    marginBottom: 15,
    color: '#666',
  },
  input: {
    height: 44,
    borderWidth: 1,
    borderColor: '#ccc',
    borderRadius: 8,
    paddingHorizontal: 10,
    marginBottom: 12,
  },
  textArea: {
    height: 100,
    textAlignVertical: 'top',
    paddingTop: 10,
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
  opportunityForm: {
    marginTop: 20,
    marginBottom: 10,
    padding: 15,
    backgroundColor: '#f9f9f9',
    borderRadius: 8,
    borderWidth: 1,
    borderColor: '#eee',
  },
  formHeading: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 15,
  },
  typeContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginBottom: 15,
  },
  typeButton: {
    paddingVertical: 8,
    paddingHorizontal: 12,
    borderRadius: 20,
    backgroundColor: '#eee',
    marginRight: 8,
    marginBottom: 8,
  },
  selectedTypeButton: {
    backgroundColor: '#007AFF',
  },
  typeText: {
    color: '#333',
  },
  selectedTypeText: {
    color: '#fff',
  },
  mediaItem: {
    marginBottom: 15,
    borderWidth: 1,
    borderColor: '#eee',
    borderRadius: 8,
    padding: 10,
  },
  mediaTypeContainer: {
    flexDirection: 'row',
    marginBottom: 10,
  },
  mediaTypeButton: {
    paddingVertical: 6,
    paddingHorizontal: 12,
    borderRadius: 14,
    backgroundColor: '#eee',
    marginRight: 8,
  },
  selectedMediaTypeButton: {
    backgroundColor: '#007AFF',
  },
  mediaTypeText: {
    fontSize: 14,
    color: '#333',
  },
  selectedMediaTypeText: {
    color: '#fff',
  },
  mediaInputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  mediaInput: {
    flex: 1,
    height: 44,
    borderWidth: 1,
    borderColor: '#ccc',
    borderRadius: 8,
    paddingHorizontal: 10,
  },
  removeButton: {
    marginLeft: 10,
    width: 30,
    height: 30,
    alignItems: 'center',
    justifyContent: 'center',
  },
  addMediaButton: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 15,
    paddingVertical: 8,
  },
  addMediaText: {
    color: '#007AFF',
    fontSize: 16,
    marginLeft: 5,
  },
});
