import React, { useState, useEffect } from 'react';
import { View, Text, TextInput, StyleSheet, ScrollView, TouchableOpacity, Alert, ActivityIndicator } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import * as SecureStore from 'expo-secure-store';
import { useRouter } from 'expo-router';
import AsyncStorage from '@react-native-async-storage/async-storage';

export default function RecruiterProfile() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [opportunity, setOpportunity] = useState(null);
  const [saving, setSaving] = useState(false);
  const router = useRouter();

  // Form state
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [location, setLocation] = useState('');
  const [type, setType] = useState('volunteer');
  const [mediaItems, setMediaItems] = useState([{ url: '', type: 'Image' }]);
  const [tags, setTags] = useState('');

  const opportunityTypes = ['volunteer', 'internship', 'job', 'event'];
  const mediaTypes = ['Image', 'Video'];

  // Add new media item
  const addMediaItem = () => {
    setMediaItems([...mediaItems, { url: '', type: 'Image' }]);
  };

  // Remove media item
  const removeMediaItem = (index) => {
    if (mediaItems.length > 1) {
      const updatedItems = [...mediaItems];
      updatedItems.splice(index, 1);
      setMediaItems(updatedItems);
    }
  };

  // Update media item
  const updateMediaItem = (index, field, value) => {
    const updatedItems = [...mediaItems];
    updatedItems[index] = { ...updatedItems[index], [field]: value };
    setMediaItems(updatedItems);
  };

  // Fetch user data on load
  useEffect(() => {
    const fetchUserData = async () => {
      try {
        const userStr = await SecureStore.getItemAsync('user');
        if (!userStr) {
          router.replace('/auth/login');
          return;
        }

        const token = await AsyncStorage.getItem('token');

        const userData = JSON.parse(userStr);
        setUser(userData);
        console.log("User data:", userData);
        // Fetch recruiter's opportunity
        fetchOpportunity(token, userData.uuid);
      } catch (error) {
        console.error('Error loading user data:', error);
        setLoading(false);
      }
    };

    fetchUserData();
  }, []);

  // Fetch recruiter's opportunity
  const fetchOpportunity = async (token, uuid) => {
    try {
      const response = await fetch(`http://192.168.1.58:8080/api/v1/opportunities/author/${uuid}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        }
      });

      // Check if response is OK before parsing
      if (!response.ok) {
        console.log(`API returned error status: ${response.status}`);
        setLoading(false);
        return;
      }

      // Get response text first to check content
      const responseText = await response.text();
      
      // Check if response is empty or not JSON
      if (!responseText || responseText.trim() === '') {
        console.log('Empty response from API');
        setLoading(false);
        return;
      }

      // Try to parse JSON
      let responseData;
      try {
        responseData = JSON.parse(responseText);
        console.log("Raw API response:", JSON.stringify(responseData));
      } catch (parseError) {
        console.error('Failed to parse response as JSON:', responseText.substring(0, 100));
        setLoading(false);
        return;
      }
      
      // Handle nested data structure - could be in data.data or just data
      let opportunityData = [];
      
      if (responseData.data && Array.isArray(responseData.data)) {
        opportunityData = responseData.data;
      } else if (responseData.data && responseData.data.data && Array.isArray(responseData.data.data)) {
        opportunityData = responseData.data.data;
      } else if (Array.isArray(responseData)) {
        opportunityData = responseData;
      }
      
      console.log("Extracted opportunity data:", opportunityData);
      
      if (opportunityData.length > 0) {
        const opp = opportunityData[0]; // Get the first opportunity
        console.log("Processing opportunity:", opp);
        
        // Extract media information
        let mediaURLs = [];
        let mediaTypes = [];
        
        if (Array.isArray(opp.media)) {
          opp.media.forEach(mediaItem => {
            if (typeof mediaItem === 'object') {
              if (mediaItem.URL) mediaURLs.push(mediaItem.URL);
              if (mediaItem.type) mediaTypes.push(mediaItem.type);
            }
          });
        }
        
        // Extract tags
        let tagNames = [];
        if (Array.isArray(opp.tags)) {
          opp.tags.forEach(tagItem => {
            if (typeof tagItem === 'object' && tagItem.tagName) {
              tagNames.push(tagItem.tagName);
            } else if (typeof tagItem === 'string') {
              tagNames.push(tagItem);
            }
          });
        }
        
        // Map the API response fields to our expected format
        const mappedOpp = {
          id: opp.uuid,
          title: opp.title || '',
          description: opp.description || '',
          location: opp.location || '',
          type: opp.opportunityType || opp.type || 'volunteer',
          author: opp.postedByUUID || opp.author,
          points: opp.points,
          approved: opp.approved,
          // Use the extracted media and tags
          tags: tagNames,
          mediaURL: mediaURLs,
          mediaType: mediaTypes
        };
        
        setOpportunity(mappedOpp);
        console.log("Mapped opportunity:", mappedOpp);
        
        // Initialize form with existing data
        setTitle(mappedOpp.title);
        setDescription(mappedOpp.description);
        setLocation(mappedOpp.location);
        setType(mappedOpp.type);
        
        // Set media items
        if (mappedOpp.mediaURL.length > 0) {
          const media = mappedOpp.mediaURL.map((url, index) => ({
            url,
            type: mappedOpp.mediaType[index] || 'Image'
          }));
          setMediaItems(media);
        }
        
        // Set tags
        if (mappedOpp.tags.length > 0) {
          setTags(mappedOpp.tags.join(', '));
        }
      } else {
        // No opportunities found or unexpected response format
        console.log('No opportunities found for this recruiter or unexpected response format');
        console.log('Response data:', responseData);
      }
      
      setLoading(false);
    } catch (error) {
      console.error('Error fetching opportunity:', error);
      setLoading(false);
    }
  };

  // Handle opportunity update
  const handleUpdateOpportunity = async () => {
    if (!title || !description || !location) {
      Alert.alert('Error', 'Please fill in all required fields');
      return;
    }

    try {
      setSaving(true);
      
      // Get token from AsyncStorage
      const token = await AsyncStorage.getItem('token');
      if (!token) {
        throw new Error('Authentication token not found. Please log in again.');
      }
      
      // Prepare tags as a simple array of strings
      const tagsArray = tags.split(',').map(tag => tag.trim()).filter(tag => tag);
      
      // Format media items as objects with type and URL
      const mediaArray = mediaItems.map(item => ({
        type: item.type,
        URL: item.url
      }));
      
      if (opportunity) {
        // Update existing opportunity - use the expected format
        const updateBody = {
          uuid: opportunity.id,
          title,
          description,
          points: opportunity.points || Math.floor(Math.random() * 500) + 100,
          location,
          opportunityType: type,
          postedByUUID: user.uuid,
          approved: opportunity.approved || false,
          tags: tagsArray,
          media: mediaArray
        };
        
        console.log("Updating opportunity with:", updateBody);
        
        const response = await fetch(`http://192.168.1.58:8080/api/v1/opportunities/`, {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(updateBody)
        });

        if (!response.ok) {
          const errorText = await response.text();
          console.error('API error response:', errorText);
          throw new Error('Failed to update opportunity');
        }

        // Try to parse the response
        try {
          const responseText = await response.text();
          console.log("Update response:", responseText);
          if (responseText) {
            const responseData = JSON.parse(responseText);
            console.log("Parsed update response:", responseData);
          }
        } catch (parseError) {
          console.error('Error parsing update response:', parseError);
        }

        Alert.alert('Success', 'Your opportunity has been updated');
      } else {
        // Create new opportunity
        const createBody = {
          title,
          description,
          points: Math.floor(Math.random() * 500) + 100,
          location,
          opportunityType: type || 'volunteer',
          postedByUUID: user.uuid,
          tags: tagsArray,
          media: mediaArray
        };
        
        console.log("Creating opportunity with:", createBody);
        
        const response = await fetch('http://192.168.1.58:8080/api/v1/opportunities', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(createBody)
        });

        if (!response.ok) {
          const errorText = await response.text();
          console.error('API error response:', errorText);
          throw new Error('Failed to create opportunity');
        }

        // Try to parse response
        try {
          const responseText = await response.text();
          console.log("Create response:", responseText);
          
          if (responseText) {
            const responseData = JSON.parse(responseText);
            console.log("Parsed create response:", responseData);
            
            let newOpp = null;
            
            // Handle different response structures
            if (responseData.data) {
              newOpp = responseData.data;
            } else if (responseData.uuid || responseData.id) {
              newOpp = responseData;
            }
            
            if (newOpp) {
              // Extract media info
              let mediaURLs = [];
              let mediaTypes = [];
              
              if (Array.isArray(newOpp.media)) {
                newOpp.media.forEach(mediaItem => {
                  if (typeof mediaItem === 'object') {
                    if (mediaItem.URL) mediaURLs.push(mediaItem.URL);
                    if (mediaItem.type) mediaTypes.push(mediaItem.type);
                  }
                });
              }
              
              // Extract tags
              let tagNames = [];
              if (Array.isArray(newOpp.tags)) {
                newOpp.tags.forEach(tagItem => {
                  if (typeof tagItem === 'object' && tagItem.tagName) {
                    tagNames.push(tagItem.tagName);
                  } else if (typeof tagItem === 'string') {
                    tagNames.push(tagItem);
                  }
                });
              }
              
              // Map to our format
              const mappedOpp = {
                id: newOpp.uuid || newOpp.id,
                title: newOpp.title,
                description: newOpp.description,
                location: newOpp.location,
                type: newOpp.opportunityType || newOpp.type,
                points: newOpp.points,
                tags: tagNames,
                mediaURL: mediaURLs,
                mediaType: mediaTypes
              };
              
              setOpportunity(mappedOpp);
              console.log("Created and mapped opportunity:", mappedOpp);
            }
          }
        } catch (parseError) {
          console.error('Failed to parse creation response:', parseError);
          // Still consider it a success even if we can't parse the response
        }
        
        Alert.alert('Success', 'Your opportunity has been created');
      }
    } catch (error) {
      console.error('Error updating opportunity:', error);
      Alert.alert('Error', error.message);
    } finally {
      setSaving(false);
    }
  };

  // Handle logout
  const handleLogout = async () => {
    try {
      await SecureStore.deleteItemAsync('user');
      await AsyncStorage.removeItem('token');
      router.replace('/auth/login');
    } catch (error) {
      console.error('Error logging out:', error);
    }
  };

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#007AFF" />
        <Text style={styles.loadingText}>Loading your profile...</Text>
      </View>
    );
  }

  return (
    <ScrollView style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Recruiter Profile</Text>
        <TouchableOpacity onPress={handleLogout} style={styles.logoutButton}>
          <Ionicons name="log-out-outline" size={24} color="#FF3B30" />
          <Text style={styles.logoutText}>Logout</Text>
        </TouchableOpacity>
      </View>

      <View style={styles.userInfo}>
        <Ionicons name="person-circle" size={60} color="#007AFF" />
        <View style={styles.userDetails}>
          <Text style={styles.username}>{user?.username || 'Recruiter'}</Text>
          <Text style={styles.role}>{user?.role || 'Recruiter'}</Text>
        </View>
      </View>

      <View style={styles.formContainer}>
        <Text style={styles.sectionTitle}>Your Opportunity</Text>
        <Text style={styles.sectionSubtitle}>Edit your opportunity details below</Text>

        <Text style={styles.label}>Title *</Text>
        <TextInput
          style={styles.input}
          value={title}
          onChangeText={setTitle}
          placeholder="Opportunity title"
        />

        <Text style={styles.label}>Description *</Text>
        <TextInput
          style={[styles.input, styles.textArea]}
          value={description}
          onChangeText={setDescription}
          placeholder="Describe your opportunity"
          multiline
          numberOfLines={6}
        />

        <Text style={styles.label}>Location *</Text>
        <TextInput
          style={styles.input}
          value={location}
          onChangeText={setLocation}
          placeholder="Where is this opportunity located?"
        />

        <Text style={styles.label}>Type</Text>
        <View style={styles.typeContainer}>
          {opportunityTypes.map((t) => (
            <TouchableOpacity
              key={t}
              style={[
                styles.typeButton,
                type === t && styles.selectedTypeButton
              ]}
              onPress={() => setType(t)}
            >
              <Text 
                style={[
                  styles.typeText,
                  type === t && styles.selectedTypeText
                ]}
              >
                {t.charAt(0).toUpperCase() + t.slice(1)}
              </Text>
            </TouchableOpacity>
          ))}
        </View>

        <Text style={styles.label}>Media</Text>
        {mediaItems.map((item, index) => (
          <View key={index} style={styles.mediaItem}>
            <View style={styles.mediaTypeContainer}>
              {mediaTypes.map(t => (
                <TouchableOpacity 
                  key={t}
                  style={[
                    styles.mediaTypeButton,
                    item.type === t && styles.selectedMediaTypeButton
                  ]}
                  onPress={() => updateMediaItem(index, 'type', t)}
                >
                  <Text 
                    style={[
                      styles.mediaTypeText,
                      item.type === t && styles.selectedMediaTypeText
                    ]}
                  >
                    {t}
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

        <Text style={styles.label}>Tags</Text>
        <TextInput
          style={styles.input}
          value={tags}
          onChangeText={setTags}
          placeholder="Enter tags separated by commas"
        />

        <TouchableOpacity 
          style={[styles.saveButton, saving && styles.savingButton]} 
          onPress={handleUpdateOpportunity}
          disabled={saving}
        >
          {saving ? (
            <ActivityIndicator size="small" color="#FFFFFF" />
          ) : (
            <Text style={styles.saveButtonText}>
              {opportunity ? 'Update Opportunity' : 'Create Opportunity'}
            </Text>
          )}
        </TouchableOpacity>
      </View>
    </ScrollView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#fff',
  },
  loadingText: {
    marginTop: 10,
    fontSize: 16,
    color: '#666',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingHorizontal: 20,
    paddingTop: 60,
    paddingBottom: 20,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  headerTitle: {
    fontSize: 22,
    fontWeight: 'bold',
  },
  logoutButton: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  logoutText: {
    color: '#FF3B30',
    marginLeft: 5,
  },
  userInfo: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 20,
  },
  userDetails: {
    marginLeft: 15,
  },
  username: {
    fontSize: 18,
    fontWeight: 'bold',
  },
  role: {
    fontSize: 16,
    color: '#666',
    marginTop: 5,
  },
  formContainer: {
    padding: 20,
  },
  sectionTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 5,
  },
  sectionSubtitle: {
    fontSize: 16,
    color: '#666',
    marginBottom: 20,
  },
  label: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 8,
    marginTop: 16,
  },
  input: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 12,
    fontSize: 16,
  },
  textArea: {
    height: 120,
    textAlignVertical: 'top',
  },
  typeContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
    marginVertical: 8,
  },
  typeButton: {
    paddingVertical: 8,
    paddingHorizontal: 16,
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
    marginTop: 5,
    marginBottom: 15,
    paddingVertical: 8,
  },
  addMediaText: {
    color: '#007AFF',
    fontSize: 16,
    marginLeft: 5,
  },
  saveButton: {
    backgroundColor: '#007AFF',
    borderRadius: 8,
    paddingVertical: 14,
    alignItems: 'center',
    marginTop: 30,
    marginBottom: 40,
  },
  savingButton: {
    backgroundColor: '#7AAEFF',
  },
  saveButtonText: {
    color: '#fff',
    fontWeight: 'bold',
    fontSize: 16,
  },
}); 