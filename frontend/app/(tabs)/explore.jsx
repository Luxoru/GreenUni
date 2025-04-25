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

const { width, height } = Dimensions.get('window');
const CARD_WIDTH = width * 0.9;
const CARD_HEIGHT = height * 0.7;

// Memoized card component for better performance
const GreenCard = memo(({ item, onLike, onDislike, expanded, setExpanded }) => {
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
    onLike(item.uuid);
  };

  const handleDislike = () => {
    collapseCard();
    onDislike(item.uuid);
  };

  return (
    <TouchableOpacity 
      activeOpacity={1}
      onPress={toggleExpand}
      style={styles.cardContainer}
    >
      <Animated.View style={[styles.card, { height: animatedHeight }]}>
        <Image
          source={{ uri: item.media[item.media.length - 1].URL }}
          style={styles.cardImage}
          resizeMode="cover"
        />
        <View style={styles.cardContent}>
          <Text style={styles.title} numberOfLines={1} ellipsizeMode="tail">
            {item.title}
          </Text>
          <Text style={styles.body} numberOfLines={expanded ? 10 : 2} ellipsizeMode="tail">
            {item.description}
          </Text>

          {expanded && (
            <View style={{ marginTop: 10 }}>
              <Text style={styles.detail}>📍 Location: {item.location || "N/A"}</Text>
              <Text style={styles.detail}>Tags: {item.tags.map(tag => tag.tagName).join(", ") || "None"}</Text>
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

const NextCardPeek = memo(({ item }) => {
  if (!item) return null;
  
  return (
    <View style={styles.peekCardContainer}>
      <View style={styles.peekCard}>
        <Image
          source={{ uri: item.imageUrl }}
          style={styles.cardImage}
          resizeMode="cover"
        />
      </View>
    </View>
  );
});

const VolunteerPage = () => {
  const [data, setData] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const [isFetchingMore, setIsFetchingMore] = useState(false);
  const [page, setPage] = useState(1);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isCardExpanded, setIsCardExpanded] = useState(false);

  const fetchData = async (from) => {
      try {
        setIsLoading(true);

        const userStr = await SecureStore.getItemAsync('user');

      if (!userStr) {
        console.warn("No user found in SecureStore.");
        return;
      }

      const user = JSON.parse(userStr);

      const response = await axios.get(`http://192.168.1.58:8080/api/v1/opportunities?from=${from}&limit=5&uuid=${user.uuid}`);

      console.log("Fetching from " + from);

      if (!response.data.data) return [];

      setPage(response.data.lastIndex);

      const enhancedData = response.data.data.map(item => ({
        ...item,
        imageUrl: item.media[item.media.length - 1]?.URL || '', 
      }));

      return enhancedData;
    } catch (error) {
      console.error(error);
      return [];
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchData(page).then((newData) => {
      setData(newData);
    });
  }, []);

  const loadMoreData = async () => {
    if (isFetchingMore) return;

    setIsFetchingMore(true);
    const newPage = page + 1;
    const moreData = await fetchData(newPage);

    if (moreData.length > 0) {
      setPage(newPage);
      setData((prevData) => [...prevData, ...moreData]);
    }

    setIsFetchingMore(false);
  };

  useEffect(() => {
    if (data.length > 0 && currentIndex >= data.length - 2) {
      loadMoreData();
    }
  }, [currentIndex, data.length]);

  const handleLike = async (id) => {
    console.log(`Liked post ${id}`);
  
    const userStr = await SecureStore.getItemAsync('user');

    if (!userStr) {
      console.warn("No user found in SecureStore.");
      return;
    }

    const user = JSON.parse(userStr);
  
    const response = await axios.post(`http://192.168.1.58:8080/api/v1/opportunities/likes/${user.uuid}/${id}`);

    console.log(response.data.success)

    if(!response.data.success){
      Alert.alert('Liking failed', `${response.data.message}`);
    }
  
    moveToNextCard();
  };
  

  const handleDislike = async (id) => {
    console.log(`Disliked post ${id}`);

    const userStr = await SecureStore.getItemAsync('user');

    if (!userStr) {
      console.warn("No user found in SecureStore.");
      return;
    }

    const user = JSON.parse(userStr);
  
    const response = await axios.post(`http://192.168.1.58:8080/api/v1/opportunities/dislikes/${user.uuid}/${id}`);

    if(response.data.success == false){
      Alert.alert('Disliking failed', `${response.data.message}`);
    }

    moveToNextCard();
  };

  const moveToNextCard = () => {
    setIsCardExpanded(false);
    setCurrentIndex(prev => prev + 1);
  };

  if (isLoading && data.length === 0) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#FF4949" />
      </View>
    );
  }

  const currentCard = data[currentIndex];
  const nextCard = data[currentIndex + 1];

  if (!currentCard && isLoading) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#FF4949" />
        <Text style={styles.loadingText}>Loading more Opportunities...</Text>
      </View>
    );
  }

  if (!currentCard && !isLoading) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#FF4949" />
        <Text style={styles.loadingText}>No more opportunities!</Text>
        <Text style={styles.loadingText}>Come back later!</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={styles.cardsContainer}>
        {!isCardExpanded && nextCard && (
          <NextCardPeek item={nextCard} />
        )}

        <GreenCard 
          item={currentCard} 
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
    width: width,
    height: CARD_HEIGHT,
    alignItems: 'center',
    justifyContent: 'center',
  },
  cardContainer: {
    width: CARD_WIDTH,
    height: CARD_HEIGHT,
    justifyContent: 'center',
    alignItems: 'center',
    position: 'absolute',
    zIndex: 2,
  },
  peekCardContainer: {
    width: CARD_WIDTH,
    height: CARD_HEIGHT,
    justifyContent: 'center',
    alignItems: 'center',
    position: 'absolute',
    zIndex: 1,
    transform: [
      { scale: 0.95 },
      { translateX: 15 }
    ],
    opacity: 0.7,
  },
  card: {
    width: '100%',
    height: '100%',
    borderRadius: 20,
    backgroundColor: '#fff',
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 5,
    },
    shadowOpacity: 0.3,
    shadowRadius: 5,
    elevation: 10,
    overflow: 'hidden',
  },
  peekCard: {
    width: '100%',
    height: '100%',
    borderRadius: 20,
    backgroundColor: '#fff',
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.2,
    shadowRadius: 3,
    elevation: 5,
    overflow: 'hidden',
  },
  cardImage: {
    width: '100%',
    height: '65%',
  },
  cardContent: {
    padding: 15,
    height: '20%',
  },
  buttonContainer: {
    height: '15%',
    flexDirection: 'row',
    justifyContent: 'space-evenly',
    alignItems: 'center',
    paddingBottom: 10,
  },
  title: {
    fontSize: 18,
    fontWeight: 'bold',
    marginBottom: 8,
    color: '#333',
  },
  body: {
    fontSize: 14,
    color: '#666',
  },
  button: {
    width: 60,
    height: 60,
    borderRadius: 30,
    justifyContent: 'center',
    alignItems: 'center',
    shadowColor: '#000',
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.25,
    shadowRadius: 3.84,
    elevation: 5,
  },
  dislikeButton: {
    backgroundColor: '#FF4949',
  },
  likeButton: {
    backgroundColor: '#4DED30',
  },
  buttonText: {
    fontSize: 24,
    color: 'white',
  },
  loadingText: {
    marginTop: 5,
    fontSize: 16,
    color: '#666',
  },
  detail: {
    fontSize: 14,
    color: '#444',
    marginBottom: 5,
  },
});

export default VolunteerPage;
