import AsyncStorage from '@react-native-async-storage/async-storage';

export const storeItem = async (key, value) => {
  try {
    await AsyncStorage.setItem(key, value);
  } catch (e) {
    console.error('Storage save error:', e);
  }
};

export const getItem = async (key) => {
  try {
    return await AsyncStorage.getItem(key);
  } catch (e) {
    console.error('Storage read error:', e);
    return null;
  }
};

export const removeItem = async (key) => {
  try {
    await AsyncStorage.removeItem(key);
  } catch (e) {
    console.error('Storage remove error:', e);
  }
};