import React, { useState, useEffect, useRef, memo } from 'react';
import { 
  View, 
  Text, 
  ActivityIndicator, 
  StyleSheet, 
  Dimensions, 
  TouchableOpacity,
  Image,
  Animated
} from 'react-native';
import axios from 'axios';

const { width, height } = Dimensions.get('window');
const CARD_WIDTH = width * 0.9;
const CARD_HEIGHT = height * 0.7;

// Memoized card component for better performance i believe
const GreenCard = memo(({ item, onLike, onDislike }) => {
  return (
    <View style={styles.cardContainer}>
      <View style={styles.card}>
        <Image
          source={{ uri: item.imageUrl }}
          style={styles.cardImage}
          resizeMode="cover"
        />
        <View style={styles.cardContent}>
          <Text style={styles.title} numberOfLines={1} ellipsizeMode="tail">
            {item.title}
          </Text>
          <Text style={styles.body} numberOfLines={2} ellipsizeMode="tail">
            {item.body}
          </Text>
        </View>
        <View style={styles.buttonContainer}>
          <TouchableOpacity 
            style={[styles.button, styles.dislikeButton]} 
            onPress={() => onDislike(item.id)}
          >
            <Text style={styles.buttonText}>✕</Text>
          </TouchableOpacity>
          <TouchableOpacity 
            style={[styles.button, styles.likeButton]} 
            onPress={() => onLike(item.id)}
          >
            <Text style={styles.buttonText}>❤</Text>
          </TouchableOpacity>
        </View>
      </View>
    </View>
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

  
  const fetchData = async (pageNum) => {
    try {
      setIsLoading(true);
      const response = await axios.get(`https://jsonplaceholder.typicode.com/posts?_page=${pageNum}&_limit=5`);
      
      
      const enhancedData = response.data.map(item => ({
        ...item,
        imageUrl: `https://picsum.photos/seed/${item.id}/400/600`
      }));
      
      return enhancedData;
    } catch (error) {
      console.error(error);
      return [];
    } finally {
      setIsLoading(false);
    }
  };

  // Load initial data
  useEffect(() => {
    fetchData(page).then((newData) => {
      setData(newData);
    });
  }, []);

  // Load more data when needed
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

  // Check if we need to load more data
  useEffect(() => {
    // If we're near the end of our data, load more
    if (data.length > 0 && currentIndex >= data.length - 2) {
      loadMoreData();
    }
  }, [currentIndex, data.length]);

  const handleLike = (id) => {
    console.log(`Liked post ${id}`);
    // Move to the next card
    moveToNextCard();
  };

  const handleDislike = (id) => {
    console.log(`Disliked post ${id}`);
    // Move to the next card
    moveToNextCard();
  };

  const moveToNextCard = () => {
    setCurrentIndex(prev => prev + 1);
  };

  if (isLoading && data.length === 0) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#FF4949" />
      </View>
    );
  }

  // Get current card and next card (for peeking)
  const currentCard = data[currentIndex];
  const nextCard = data[currentIndex + 1];

  if (!currentCard) {
    return (
      <View style={styles.container}>
        <ActivityIndicator size="large" color="#FF4949" />
        <Text style={styles.loadingText}>Loading more Opportunities...</Text>
      </View>
    );
  }

  return (
    <View style={styles.container}>
      <View style={styles.cardsContainer}>
        {/* Next card peeking */}
        {nextCard && (
          <NextCardPeek item={nextCard} />
        )}
        
        {/* Current card*/}
        <GreenCard 
          item={currentCard} 
          onLike={handleLike} 
          onDislike={handleDislike}
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
    marginTop: 15,
    fontSize: 16,
    color: '#666',
  }
});

export default VolunteerPage;