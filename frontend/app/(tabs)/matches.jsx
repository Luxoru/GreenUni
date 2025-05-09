import React, { useState, useEffect } from 'react';
import { 
  View, 
  Text, 
  StyleSheet, 
  FlatList, 
  Image, 
  TouchableOpacity, 
  ActivityIndicator,
  Alert 
} from 'react-native';
import { useRouter } from 'expo-router';
import * as SecureStore from 'expo-secure-store';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { Ionicons } from '@expo/vector-icons';
import config from '../../utils/config';
export default function MatchesScreen() {
  const [loading, setLoading] = useState(true);
  const [matches, setMatches] = useState([]);
  const [user, setUser] = useState(null);
  const router = useRouter();

  // Load user data and fetch matches
  useEffect(() => {
    const loadUserAndMatches = async () => {
      try {
        // Get user data from secure storage
        const userStr = await SecureStore.getItemAsync('user');
        if (!userStr) {
          console.warn("No user found in SecureStore");
          router.replace('/auth/login');
          return;
        }

        const userData = JSON.parse(userStr);
        setUser(userData);

        // Fetch matches for this user based on role
        if (userData.role === 'Recruiter') {
          await fetchRecruiterMatches(userData.uuid);
        } else {
          await fetchStudentMatches(userData.uuid);
        }
      } catch (error) {
        console.error("Error loading user data:", error);
        setLoading(false);
      }
    };

    loadUserAndMatches();
  }, []);

  // Fetch recruiter's opportunities first, then get likes for those opportunities
  const fetchRecruiterMatches = async (recruiterId) => {
    try {
      setLoading(true);
      
      // Get authentication token
      const token = await AsyncStorage.getItem('token');
      if (!token) {
        Alert.alert("Session Expired", "Please login again");
        router.replace('/auth/login');
        return;
      }

      // API call to get matches
      const response = await fetch(`${config.apiURL}/api/v1/match/${recruiterId}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
      }

      const data = await response.json();
      console.log("Recruiter matches data:", data);
      
      if (data.success && data.data) {
        
        const studentsWithMedia = await Promise.all(
          data.data.map(async (student) => {
            try {
              const studentResponse = await fetch(`${config.apiURL}/api/v1/student/${student.uuid}`, {
                method: 'GET',
                headers: {
                  'Content-Type': 'application/json'
                }
              });
              
              if (!studentResponse.ok) {
                console.log("Student response not ok:", studentResponse);
                return student;
              }
              
              const studentData = await studentResponse.json();

              console.log("Student data:", studentData);
              
              return {
                ...student,

                profilePic: studentData.data.profilePic
              };

            } catch (error) {
              console.error("Error fetching student details:", error);
              return student;
            }
          })
        );
        
        setMatches(studentsWithMedia);
      } else {
        setMatches([]);
      }
    } catch (error) {
      console.error("Error fetching recruiter matches:", error);
      Alert.alert("Error", "Failed to load matches. Please try again later.");
    } finally {
      setLoading(false);
    }
  };

  // Fetch student matches directly from the match endpoint
  const fetchStudentMatches = async (studentId) => {
    try {
      setLoading(true);
      
      // Get authentication token
      const token = await AsyncStorage.getItem('token');
      if (!token) {
        Alert.alert("Session Expired", "Please login again");
        router.replace('/auth/login');
        return;
      }

      // API call to get matches
      const response = await fetch(`${config.apiURL}/api/v1/match/${studentId}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json'
        }
      });

      if (!response.ok) {
        throw new Error(`API error: ${response.status}`);
      }

      const data = await response.json();
      console.log("Student matches data:", data);
      
      if (data.success && data.data) {
        const recruitersWithMedia = await Promise.all(
          data.data.map(async (recruiter) => {
            try {
              const oppResponse = await fetch(`${config.apiURL}/api/v1/opportunities/author/${recruiter.uuid}`, {
                method: 'GET',
                headers: {
                  'Authorization': `Bearer ${token}`,
                  'Content-Type': 'application/json'
                }
              });
              
              if (!oppResponse.ok) {
                return recruiter; 
              }
              
              const oppData = await oppResponse.json();
              
              
              if (oppData.data && oppData.data.length > 0 && 
                  oppData.data[0].media && oppData.data[0].media.length > 0) {
                
                
                const mediaItem = oppData.data[0].media.find(item => item.URL);
                if (mediaItem) {
                  return {
                    ...recruiter,
                    profilePic: mediaItem.URL
                  };
                }
              }
              return recruiter;
            } catch (error) {
              console.error("Error fetching recruiter opportunity:", error);
              return recruiter;
            }
          })
        );
        
        setMatches(recruitersWithMedia);
      } else {
        setMatches([]); //Resets matches to 0
      }
    } catch (error) {
      console.error("Error fetching student matches:", error);
      Alert.alert("Error", "Failed to load matches. Please try again later.");
    } finally {
      setLoading(false);
    }
  };


  const handleRefresh = () => {
    if (user) {
      if (user.role === 'Recruiter') {
        fetchRecruiterMatches(user.uuid);
      } else {
        fetchStudentMatches(user.uuid);
      }
    }
  };


  const renderMatchItem = ({ item }) => {

    if (user?.role === 'Recruiter') {
      console.log("Recruiter match item:", item);
      return (
        <TouchableOpacity 
          style={styles.matchCard}
          onPress={() => Alert.alert(
            'Student Details', 
            `${item.username}\n\nEmail: ${item.email}\n\n${item.description || ''}`
          )}
        >
          <View style={styles.avatarContainer}>
            {item.profilePic ? (
              <Image 
                source={{ uri: item.profilePic }} 
                style={styles.matchAvatar} 
                resizeMode="cover"
              />
            ) : (
              <Ionicons name="person-circle" size={60} color="#007AFF" />
            )}
          </View>
          <View style={styles.matchDetails}>
            <Text style={styles.matchName}>{item.username}</Text>
            <Text style={styles.matchEmail}>{item.email}</Text>
            {item.description && (
              <Text style={styles.matchDescription} numberOfLines={2}>
                {item.description}
              </Text>
            )}
            
            <View style={styles.actionButtons}>
              <TouchableOpacity style={styles.messageButton}>
                <Ionicons name="mail-outline" size={20} color="#FFFFFF" />
                <Text style={styles.buttonText}>Contact</Text>
              </TouchableOpacity>
            </View>
          </View>
        </TouchableOpacity>
      );
    } 
    else {
      console.log("Student match item:", item);
      return (
        <TouchableOpacity 
          style={styles.matchCard}
          onPress={() => Alert.alert(
            'Recruiter Details', 
            `Email: ${item.email}\n\nRole: ${item.role || 'Recruiter'}`
          )}
        >
          <View style={styles.avatarContainer}>
            {item.profilePic ? (
              <Image 
                source={{ uri: item.profilePic }} 
                style={styles.matchAvatar} 
                resizeMode="cover"
              />
            ) : (
              <Ionicons name="person-circle" size={60} color="#007AFF" />
            )}
          </View>
          <View style={styles.matchDetails}>
            <Text style={styles.matchName}>{item.username}</Text>
            <Text style={styles.matchEmail}>{item.email}</Text>
            <Text style={styles.matchRole}>{item.role || 'Recruiter'}</Text>
            
            <View style={styles.actionButtons}>
              <TouchableOpacity style={styles.messageButton}>
                <Ionicons name="mail-outline" size={20} color="#FFFFFF" />
                <Text style={styles.buttonText}>Contact</Text>
              </TouchableOpacity>
            </View>
          </View>
        </TouchableOpacity>
      );
    }
  };

  const EmptyMatches = () => (
    <View style={styles.emptyContainer}>
      <Ionicons name="heart-outline" size={80} color="#CCCCCC" />
      <Text style={styles.emptyTitle}>No Matches Yet</Text>
      <Text style={styles.emptyText}>
        {user?.role === 'Recruiter' 
          ? "No students have liked your opportunity yet. Check back in a bit lad!" 
          : "When you like opportunities and recruiters like you back, they'll appear here! Don't worry, you'll get they'll come flying"}
      </Text>
      <TouchableOpacity 
        style={styles.exploreButton}
        onPress={() => router.push(user?.role === 'Recruiter' ? '/recruiterExplore' : '/explore')}
      >
        <Text style={styles.exploreButtonText}>Keep Exploring</Text>
      </TouchableOpacity>
    </View>
  );

  if (loading) {
    return (
      <View style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#007AFF" />
        <Text style={styles.loadingText}>Loading your matches...</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={styles.header}>
        <Text style={styles.headerTitle}>Your Matches</Text>
        <TouchableOpacity onPress={handleRefresh} style={styles.refreshButton}>
          <Ionicons name="refresh" size={24} color="#007AFF" />
        </TouchableOpacity>
      </View>

      <FlatList
        data={matches}
        renderItem={renderMatchItem}
        keyExtractor={(item) => item.studentID || item.uuid}
        contentContainerStyle={matches.length === 0 ? {flex: 1} : {paddingBottom: 20}}
        ListEmptyComponent={EmptyMatches}
      />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#F8F8F8',
  },
  loadingContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#F8F8F8',
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
    backgroundColor: '#FFFFFF',
    borderBottomWidth: 1,
    borderBottomColor: '#EEEEEE',
  },
  headerTitle: {
    fontSize: 22,
    fontWeight: 'bold',
  },
  refreshButton: {
    padding: 10,
  },
  matchCard: {
    flexDirection: 'row',
    backgroundColor: '#FFFFFF',
    borderRadius: 12,
    marginHorizontal: 16,
    marginTop: 16,
    padding: 16,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.1,
    shadowRadius: 4,
    elevation: 2,
  },
  avatarContainer: {
    marginRight: 16,
    justifyContent: 'center',
  },
  matchAvatar: {
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: '#E1E1E1',
  },
  matchDetails: {
    flex: 1,
  },
  matchName: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 4,
  },
  matchRole: {
    fontSize: 14,
    color: '#666666',
    marginBottom: 12,
  },
  matchEmail: {
    fontSize: 14,
    color: '#666666',
    marginBottom: 4,
  },
  matchDescription: {
    fontSize: 14,
    color: '#333333',
    marginBottom: 12,
  },
  actionButtons: {
    flexDirection: 'row',
    marginTop: 8,
  },
  messageButton: {
    flexDirection: 'row',
    alignItems: 'center',
    backgroundColor: '#007AFF',
    paddingVertical: 8,
    paddingHorizontal: 12,
    borderRadius: 20,
  },
  buttonText: {
    color: '#FFFFFF',
    fontWeight: '600',
    marginLeft: 6,
  },
  emptyContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 32,
  },
  emptyTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#333333',
    marginTop: 16,
    marginBottom: 8,
  },
  emptyText: {
    fontSize: 16,
    color: '#666666',
    textAlign: 'center',
    marginBottom: 24,
  },
  exploreButton: {
    backgroundColor: '#007AFF',
    paddingVertical: 12,
    paddingHorizontal: 24,
    borderRadius: 24,
  },
  exploreButtonText: {
    color: '#FFFFFF',
    fontWeight: 'bold',
    fontSize: 16,
  },
}); 