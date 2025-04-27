import React, { useState, useRef, useEffect } from 'react';
import { 
  View, 
  Text, 
  TextInput, 
  Button, 
  Image, 
  TouchableOpacity, 
  ScrollView, 
  StyleSheet,
  KeyboardAvoidingView,
  Platform,
  SafeAreaView,
  ActivityIndicator
} from 'react-native';
import { useRouter } from 'expo-router';
import * as ImagePicker from 'expo-image-picker';
import AsyncStorage from '@react-native-async-storage/async-storage';
import config from '../../utils/config';

export default function ProfilePage() {
  const [loading, setLoading] = useState(true);
  const [profilePic, setProfilePic] = useState(null);
  const [description, setDescription] = useState('');
  const [favoriteTags, setFavoriteTags] = useState([]);
  const [dislikedTags, setDislikedTags] = useState([]);
  const [newFavTag, setNewFavTag] = useState('');
  const [newDislikeTag, setNewDislikeTag] = useState('');
  const [studentID, setStudentID] = useState(null);
  const [studentEmail, setStudentEmail] = useState('');
  const router = useRouter();
  
  // Using refs for both TextInputs to maintain focus  -- TODO: Make this a hook
  const favTagInputRef = useRef(null);
  const dislikeTagInputRef = useRef(null);

  // Fetch user profile data on component mount -- Dont delete pls whatever u do :)
  useEffect(() => {
    fetchProfileData();
  }, []);

  const getAuthToken = async () => {
    try {
      const token = await AsyncStorage.getItem('token');
      return token;
    } catch (error) {
      console.error('Error getting auth token:', error);
      return null;
    }
  };

  const fetchProfileData = async () => {
    try {
      setLoading(true);
      const token = await getAuthToken();
      
      if (!token) {
        console.error('No authentication token found');
        return;
      }

      const response = await fetch(`${config.apiURL}/api/v1/student/me`, {
        method: "GET",
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });
      const result = await response.json();
      
      if (result.success) {
        const { data } = result;
        setStudentID(data.studentID);
        setStudentEmail(data.studentEmail || '');
        setDescription(data.description || '');
        setProfilePic(data.profilePic || null);
        setFavoriteTags(data.tagsLiked || []);
        setDislikedTags(data.tagsDisliked || []);
      } else {
        console.error('Failed to load profile data');
      }
    } catch (error) {
      console.error('Error fetching profile data:', error);
    } finally {
      setLoading(false);
    }
  };

  const pickImage = async () => {
    let result = await ImagePicker.launchImageLibraryAsync({
      mediaTypes: ImagePicker.MediaTypeOptions.Images, //FIX THE DEPRECATION SHAF
      allowsEditing: true,
      aspect: [1, 1],
      quality: 0.5,
    });

    if (!result.canceled) {
      const uploadedUrl = await uploadProfilePic(result.assets[0].uri);
      if (uploadedUrl) {
        setProfilePic(uploadedUrl);
      }
    }
  };

  const addFavoriteTag = () => {
    if (newFavTag.trim() !== '' && !favoriteTags.includes(newFavTag.trim())) {
      setFavoriteTags([...favoriteTags, newFavTag.trim()]);
      setNewFavTag('');

      if (favTagInputRef.current) {
        favTagInputRef.current.focus();
      }
    }
  };

  const addDislikedTag = () => {
    if (newDislikeTag.trim() !== '' && !dislikedTags.includes(newDislikeTag.trim())) {
      setDislikedTags([...dislikedTags, newDislikeTag.trim()]);
      setNewDislikeTag('');

      if (dislikeTagInputRef.current) {
        dislikeTagInputRef.current.focus();
      }
    }
  };

  const removeTag = (tag, list, setList) => {
    setList(list.filter(t => t !== tag));
  };

  const handleSave = async () => {
    try {
      setLoading(true);
      const token = await getAuthToken();
      
      if (!token) {
        alert('Authentication failed. Please log in again.');
        return;
      }
      
      const profileData = {
        studentID,
        studentEmail,
        description: description || '',
        profilePic: profilePic || '',
        tagsLiked: favoriteTags || [],
        tagsDisliked: dislikedTags || [],
      };
      
      const response = await fetch(`${config.apiURL}/api/v1/student/me`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(profileData),
      });
      
      const result = await response.json();
      
      if (result.success) {
        alert('Profile updated successfully!');
      } else {
        alert('Failed to update profile. Please try again.');
      }
    } catch (error) {
      console.error('Error updating profile:', error);
      alert('An error occurred while updating your profile.');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      await AsyncStorage.removeItem('token');
      console.log('Logged out successfully');
      router.replace("/auth/login")
    } catch (error) {
      console.error('Error during logout:', error);
    }
  };

  if (loading) {
    return (
      <SafeAreaView style={[styles.container, styles.loadingContainer]}>
        <ActivityIndicator size="large" color="#3498db" />
        <Text style={styles.loadingText}>Loading profile...</Text>
      </SafeAreaView>
    );
  }

  //Specific SAVs for keyboard avoidance
  return (
    <SafeAreaView style={styles.container}>
      <KeyboardAvoidingView
        behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
        style={styles.keyboardAvoid}
      >
        <ScrollView contentContainerStyle={styles.scrollContent}>
          <Text style={styles.title}>Edit Profile</Text>
          
          <View style={styles.profileImageContainer}>
            {profilePic ? (
              <Image
                source={{ uri: profilePic }}
                style={styles.profileImage}
              />
            ) : (
              <View style={styles.placeholderImage}>
                <Text style={styles.placeholderText}>No Image</Text>
              </View>
            )}
            <TouchableOpacity 
              onPress={pickImage} 
              style={styles.changePhotoButton}
            >
              <Text style={styles.changePhotoText}>Change Profile Picture</Text>
            </TouchableOpacity>
          </View>

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Email</Text>
            <TextInput
              value={studentEmail}
              onChangeText={setStudentEmail}
              placeholder="Your email address"
              style={styles.emailInput}
              keyboardType="email-address"
              autoCapitalize="none"
            />
          </View>

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Description</Text>
            <TextInput
              value={description}
              onChangeText={setDescription}
              placeholder="Tell us about yourself..."
              style={styles.descriptionInput}
              multiline
            />
          </View>

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Favorite Tags</Text>
            <View style={styles.tagInputContainer}>
              <TextInput
                ref={favTagInputRef}
                value={newFavTag}
                onChangeText={setNewFavTag}
                placeholder="Enter a favorite tag"
                style={styles.tagInput}
                returnKeyType="done"
                onSubmitEditing={addFavoriteTag}
              />
              <TouchableOpacity 
                style={styles.addButton}
                onPress={addFavoriteTag}
              >
                <Text style={styles.addButtonText}>Add</Text>
              </TouchableOpacity>
            </View>
            <View style={styles.tagsContainer}>
              {favoriteTags.map(tag => (
                <TouchableOpacity
                  key={`fav-${tag}`}
                  onPress={() => removeTag(tag, favoriteTags, setFavoriteTags)}
                  style={[styles.tag, styles.favoriteTag]}
                >
                  <Text style={styles.tagText}>{tag}</Text>
                  <Text style={styles.removeIcon}> ✖</Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          <View style={styles.section}>
            <Text style={styles.sectionTitle}>Disliked Tags</Text>
            <View style={styles.tagInputContainer}>
              <TextInput
                ref={dislikeTagInputRef}
                value={newDislikeTag}
                onChangeText={setNewDislikeTag}
                placeholder="Enter a disliked tag"
                style={styles.tagInput}
                returnKeyType="done"
                onSubmitEditing={addDislikedTag}
              />
              <TouchableOpacity 
                style={[styles.addButton, styles.dislikeButton]}
                onPress={addDislikedTag}
              >
                <Text style={styles.addButtonText}>Add</Text>
              </TouchableOpacity>
            </View>
            <View style={styles.tagsContainer}>
              {dislikedTags.map(tag => (
                <TouchableOpacity
                  key={`dislike-${tag}`}
                  onPress={() => removeTag(tag, dislikedTags, setDislikedTags)}
                  style={[styles.tag, styles.dislikedTag]}
                >
                  <Text style={styles.tagText}>{tag}</Text>
                  <Text style={styles.removeIcon}> ✖</Text>
                </TouchableOpacity>
              ))}
            </View>
          </View>

          <TouchableOpacity 
            style={styles.saveButton}
            onPress={handleSave}
            disabled={loading}
          >
            {loading ? (
              <ActivityIndicator color="#fff" size="small" />
            ) : (
              <Text style={styles.saveButtonText}>Save Changes</Text>
            )}
          </TouchableOpacity>
          
          <TouchableOpacity 
            style={styles.logoutButton} 
            onPress={handleLogout}
            disabled={loading}
          >
            <Text style={styles.logoutButtonText}>Logout</Text>
          </TouchableOpacity>
        </ScrollView>
      </KeyboardAvoidingView>
    </SafeAreaView>
  );
}

const uploadProfilePic = async (imageUri) => {
  const formData = new FormData();
  const token = await AsyncStorage.getItem('userToken');

  formData.append('file', {
    uri: imageUri,
    type: 'image/jpeg',
    name: 'profile.jpg', //Only using jpeg for now
  });
  formData.append('upload_preset', 'profile-upload'); 
  
  try {
    const res = await fetch('https://api.cloudinary.com/v1_1/dnz0ljksa/image/upload', {
      method: 'POST',
      body: formData,
    });

    const data = await res.json();
    return data.secure_url;
  } catch (err) {
    console.error('Upload failed', err);
    return null;
  }
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#f5f5f5',
  },
  loadingContainer: {
    justifyContent: 'center',
    alignItems: 'center',
  },
  loadingText: {
    marginTop: 10,
    fontSize: 16,
    color: '#333',
  },
  keyboardAvoid: {
    flex: 1,
  },
  scrollContent: {
    padding: 20,
    paddingBottom: 40,
  },
  title: {
    fontSize: 24,
    fontWeight: 'bold',
    textAlign: 'center',
    marginBottom: 20,
    color: '#333',
  },
  profileImageContainer: {
    alignItems: 'center',
    marginBottom: 24,
    paddingTop: 10,
  },
  profileImage: {
    width: 120,
    height: 120,
    borderRadius: 60,
    borderWidth: 3,
    borderColor: '#fff',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.2,
    shadowRadius: 4,
    elevation: 4,
  },
  placeholderImage: {
    width: 120,
    height: 120,
    borderRadius: 60,
    backgroundColor: '#ddd',
    alignItems: 'center',
    justifyContent: 'center',
    borderWidth: 3,
    borderColor: '#fff',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.2,
    shadowRadius: 4,
    elevation: 4,
  },
  placeholderText: {
    color: '#888',
  },
  changePhotoButton: {
    marginTop: 12,
    padding: 8,
  },
  changePhotoText: {
    color: '#3498db',
    fontWeight: '500',
  },
  section: {
    marginBottom: 24,
  },
  sectionTitle: {
    fontSize: 18,
    fontWeight: '600',
    marginBottom: 10,
    color: '#333',
  },
  emailInput: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 12,
    backgroundColor: '#fff',
  },
  descriptionInput: {
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 12,
    backgroundColor: '#fff',
    minHeight: 100,
    textAlignVertical: 'top',
  },
  tagInputContainer: {
    flexDirection: 'row',
    marginBottom: 12,
  },
  tagInput: {
    flex: 1,
    borderWidth: 1,
    borderColor: '#ddd',
    borderRadius: 8,
    padding: 12,
    backgroundColor: '#fff',
  },
  addButton: {
    backgroundColor: '#3498db',
    paddingHorizontal: 16,
    justifyContent: 'center',
    alignItems: 'center',
    borderRadius: 8,
    marginLeft: 8,
  },
  dislikeButton: {
    backgroundColor: '#e74c3c',
  },
  addButtonText: {
    color: '#fff',
    fontWeight: '600',
  },
  tagsContainer: {
    flexDirection: 'row',
    flexWrap: 'wrap',
  },
  tag: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 8,
    margin: 4,
    borderRadius: 20,
  },
  favoriteTag: {
    backgroundColor: '#2ecc71',
  },
  dislikedTag: {
    backgroundColor: '#e74c3c',
  },
  tagText: {
    color: '#fff',
    fontWeight: '500',
  },
  removeIcon: {
    color: '#fff',
    marginLeft: 4,
  },
  saveButton: {
    backgroundColor: '#3498db',
    padding: 16,
    borderRadius: 8,
    alignItems: 'center',
    marginTop: 8,
  },
  saveButtonText: {
    color: '#fff',
    fontWeight: '600',
    fontSize: 16,
  },
  logoutButton: {
    backgroundColor: '#e74c3c',
    padding: 16,
    borderRadius: 8,
    alignItems: 'center',
    marginTop: 16,
  },
  logoutButtonText: {
    color: '#fff',
    fontWeight: '600',
    fontSize: 16,
  },
});