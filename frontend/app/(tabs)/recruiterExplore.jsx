import React, { useState, useEffect, useRef, memo } from 'react';
import { 
  View, 
  Text, 
  ActivityIndicator, 
  StyleSheet, 
  Dimensions, 
  TouchableOpacity,
  Image,
  Animated,
  Alert
} from 'react-native';
import axios from 'axios';
import * as SecureStore from 'expo-secure-store';
import AsyncStorage from '@react-native-async-storage/async-storage';
import config from '../../utils/config';
const { width, height } = Dimensions.get('window');
const CARD_WIDTH = width * 0.9;
const CARD_HEIGHT = height * 0.75; 

// Memoized Student Card bettter efficiency liek that
const StudentCard = memo(({ student, onLike, onDislike, expanded, setExpanded }) => {
  const animatedHeight = useRef(new Animated.Value(CARD_HEIGHT)).current;

  const toggleExpand = () => {
    const toValue = expanded ? CARD_HEIGHT : height * 0.9;

    Animated.spring(animatedHeight, {
      toValue,
      useNativeDriver: false,
    }).start();

    setExpanded(!expanded);
  };

  const collapseCard = () => {
    Animated.spring(animatedHeight, {
      toValue: CARD_HEIGHT,
      useNativeDriver: false,
    }).start();
    setExpanded(false);
  };

  const handleLike = () => {
    collapseCard();
    onLike(student.studentID);
  };

  const handleDislike = () => {
    collapseCard();
    onDislike(student.studentID);
  };

  return (
    <TouchableOpacity 
      activeOpacity={1}
      onPress={toggleExpand}
      style={styles.cardContainer}
    >
      <Animated.View style={[styles.card, { height: animatedHeight }]}>
        <Image
          source={{ uri: student.profilePic || 'https://via.placeholder.com/300' }}
          style={styles.cardImage}
          resizeMode="cover"
        />
        <View style={styles.cardContent}>
          <Text style={styles.title} numberOfLines={1} ellipsizeMode="tail">
            {student.studentName}
          </Text>
          <Text style={styles.body} numberOfLines={expanded ? 10 : 2} ellipsizeMode="tail">
            {student.description || 'No bio available.'}
          </Text>

          {expanded && (
            <View style={{ marginTop: 10 }}>
              <Text style={styles.detail}>📍 Email: {student.studentEmail || "N/A"}</Text>

            </View>
          )}
        </View>
        <View style={styles.buttonContainer}>
          <TouchableOpacity 
            style={[styles.button, styles.dislikeButton]} 
            onPress={handleDislike}
          >
            <Text style={styles.buttonText}>✕</Text>
          </TouchableOpacity>
          <TouchableOpacity 
            style={[styles.button, styles.likeButton]} 
            onPress={handleLike}
          >
            <Text style={styles.buttonText}>❤</Text>
          </TouchableOpacity>
        </View>
      </Animated.View>
    </TouchableOpacity>
  );
});

const RecruiterPage = () => {
  const [students, setStudents] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isFetchingMore, setIsFetchingMore] = useState(false);
  const [page, setPage] = useState(0);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isCardExpanded, setIsCardExpanded] = useState(false);
  const [opportunityUUID, setOpportunityUUID] = useState(null);
  const [initialLoading, setInitialLoading] = useState(true);

  // Fetch the recruiter's opportunity UUID
  useEffect(() => {
    const fetchOpportunityUUID = async () => {
      try {
        const userStr = await SecureStore.getItemAsync('user');
        if (!userStr) {
          console.warn("No user found in SecureStore.");
          setInitialLoading(false);
          return;
        }
        
        const user = JSON.parse(userStr);
        console.log("Fetching opportunities for user:", user.uuid);
        
       
        const token = await AsyncStorage.getItem('token');
        if (!token) {
          console.warn("No token found in AsyncStorage.");
          setInitialLoading(false);
          return;
        }
        
        //Why are u using axios bruh
        const response = await axios.get(
          `${config.apiURL}/api/v1/opportunities/author/${user.uuid}`,
          {
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`
            }
          }
        );
        
        let opportunityData = null;
        
        if (response.data.success && response.data.data && response.data.data.length > 0) {
          opportunityData = response.data.data[0];
        } else if (response.data.data && Array.isArray(response.data.data) && response.data.data.length > 0) {
          opportunityData = response.data.data[0];
        } else if (Array.isArray(response.data) && response.data.length > 0) {
          opportunityData = response.data[0];
        }
        if (opportunityData) {
          const uuid = opportunityData.uuid;
          console.log("Found recruiter opportunity UUID:", uuid);
          setOpportunityUUID(uuid);
          fetchStudents(0, uuid).then((newStudents) => {
            setStudents(newStudents); // DONT DELETE THIS. IF DELETED NOT SHOWN
          });
        } else {
          console.warn("No opportunities found for this recruiter");
          setInitialLoading(false);
        }
      } catch (error) {
        console.error("Error fetching opportunity UUID:", error);
        setInitialLoading(false);
      }
    };
    
    fetchOpportunityUUID();
  }, []);

  const fetchStudents = async (from, uuid) => {
    try {
      setIsLoading(true);
      

      const fromIndex = from !== undefined ? from : 0;

      const opportunityId = uuid || opportunityUUID;
      
      if (!opportunityId) {
        console.warn("No opportunity UUID available for fetching students");
        return [];
      }
      

      const response = await axios.get(`${config.apiURL}/api/v1/opportunities/likes/${opportunityId}?from=${fromIndex}&limit=5`);

      
      if (!response.data.data || !response.data.data.likes) {
        setInitialLoading(false);
        return [];
      }
      
      setPage(response.data.data.lastIndex);
      setInitialLoading(false);
      

      return response.data.data.likes.map(likeData => {

        console.log("Like data structure:", JSON.stringify(likeData));

        if (likeData.studentID || likeData.uuid) {
          return likeData;
        }
        

        if (likeData.student) {
          return likeData.student;
        }
        

        if (Array.isArray(likeData) && likeData.length > 0) {
          return likeData[0];
        }
        

        return likeData;
      });
    } catch (error) {
      console.error("Error fetching students:", error);
      setInitialLoading(false);
      return [];
    } finally {
      setIsLoading(false);
    }
  };
  

  useEffect(() => {

    if (opportunityUUID && students.length === 0 && !isLoading) {
      fetchStudents(0).then((newStudents) => {
        setStudents(newStudents);
      });
    }
  }, [opportunityUUID]);

  const loadMoreStudents = async () => {
    if (isFetchingMore || !opportunityUUID) return;

    setIsFetchingMore(true);
    const newPage = page + 1;
    const moreStudents = await fetchStudents(newPage);
    
    if (moreStudents.length > 0) {
      setPage(newPage);
      setStudents((prevStudents) => [...prevStudents, ...moreStudents]);
    }

    setIsFetchingMore(false);
  };

  useEffect(() => {
    if (students.length > 0 && currentIndex >= students.length - 2) {
      loadMoreStudents();
    }
  }, [currentIndex, students.length]);

  const handleLike = async (studentId) => {
    console.log(`Liked student ${studentId}`);

    const recruiterStr = await SecureStore.getItemAsync('user');
    if (!recruiterStr) return;

    const recruiter = JSON.parse(recruiterStr);

    const response = await axios.post(`${config.apiURL}/api/v1/match?uuid1=${recruiter.uuid}&uuid2=${studentId}`);
    if (!response.data.success) {
      Alert.alert('Matching failed', response.data.message);
    }

    moveToNextCard();
  };

  const handleDislike = async (studentId) => {
    console.log(`Disliked student ${studentId}`);
    //No endpoint rn idk what to do
    moveToNextCard();
  };

  const moveToNextCard = () => {
    setIsCardExpanded(false);
    setCurrentIndex(prev => prev + 1);
  };

  if (initialLoading || (isLoading && students.length === 0)) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#FF4949" />
        <Text style={styles.loadingText}>Loading...</Text>
      </View>
    );
  }

  if (!opportunityUUID) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#FF4949" />
        <Text style={styles.loadingText}>No more students!</Text>
        <Text style={styles.loadingText}>Check back later!</Text>
      </View>
    );
  }

  const currentStudent = students[currentIndex];

  if (!currentStudent && isLoading) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#FF4949" />
        <Text style={styles.loadingText}>Loading more students...</Text>
      </View>
    );
  }

  if (!currentStudent && !isLoading) {
    return (
      <View style={styles.container}>
        <Text style={styles.loadingText}>No more students!</Text>
        <Text style={styles.loadingText}>Check back later!</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={styles.cardsContainer}>
        <StudentCard 
          student={currentStudent} 
          onLike={handleLike} 
          onDislike={handleDislike}
          expanded={isCardExpanded}
          setExpanded={setIsCardExpanded}
        />
      </View>
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#f5f5f5',
  },
  cardsContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
  },
  cardContainer: {
    width: CARD_WIDTH,
    justifyContent: 'center',
    alignItems: 'center',
  },
  card: {
    width: '100%',
    backgroundColor: 'white',
    borderRadius: 10,
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 2 },
    shadowOpacity: 0.25,
    shadowRadius: 3.84,
    elevation: 5,
    overflow: 'hidden',
  },
  cardImage: {
    width: '100%',
    height: '60%',
  },
  cardContent: {
    padding: 15,
    flex: 1,
  },
  title: {
    fontSize: 20,
    fontWeight: 'bold',
    marginBottom: 5,
  },
  body: {
    fontSize: 16,
    color: '#333',
  },
  detail: {
    fontSize: 14,
    color: '#555',
    marginBottom: 5,
  },
  buttonContainer: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    padding: 15,
  },
  button: {
    width: 60,
    height: 60,
    borderRadius: 30,
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: { width: 0, height: 1 },
    shadowOpacity: 0.22,
    shadowRadius: 2.22,
    elevation: 3,
  },
  likeButton: {
    backgroundColor: '#FF4949',
  },
  dislikeButton: {
    backgroundColor: '#CCC',
  },
  buttonText: {
    fontSize: 24,
    color: 'white',
  },
  loadingText: {
    marginTop: 10,
    fontSize: 16,
    color: '#666',
  },
  errorText: {
    fontSize: 18,
    fontWeight: 'bold',
    color: '#FF4949',
    marginBottom: 10,
  },
  subText: {
    fontSize: 16,
    color: '#666',
  },
});

export default RecruiterPage;
