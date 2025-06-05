export const API_URL = 'http://83.217.210.25:8080/api';

const getAuthHeader = () => {
  const token = localStorage.getItem('token');
  return {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  };
};

const fetchWithCredentials = async (url, options = {}) => {
  const defaultOptions = {
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json',
      ...options.headers
    }
  };

  return fetch(url, { ...defaultOptions, ...options });
};

export const getInboxMessages = async () => {
  try {
    const response = await fetchWithCredentials(`${API_URL}/messages/inbox`, {
      headers: getAuthHeader()
    });
    
    if (!response.ok) {
      throw new Error('Не удалось получить входящие сообщения');
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error fetching inbox messages:', error);
    return [];
  }
};

export const getSentMessages = async () => {
  try {
    const response = await fetch(`${API_URL}/messages/sent`, {
      headers: getAuthHeader()
    });
    
    if (!response.ok) {
      throw new Error('Не удалось получить отправленные сообщения');
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error fetching sent messages:', error);
    return [];
  }
};

export const getSpamMessages = async () => {
  try {
    const response = await fetch(`${API_URL}/messages/spam`, {
      headers: getAuthHeader()
    });
    
    if (!response.ok) {
      throw new Error('Не удалось получить сообщения из спама');
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error fetching spam messages:', error);
    return [];
  }
};

export const getTrashMessages = async () => {
  try {
    const response = await fetch(`${API_URL}/messages/trash`, {
      headers: getAuthHeader()
    });
    
    if (!response.ok) {
      throw new Error('Не удалось получить удаленные сообщения');
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error fetching trash messages:', error);
    return [];
  }
};

export const getMessage = async (id) => {
  try {
    const response = await fetch(`${API_URL}/messages/${id}`, {
      headers: getAuthHeader()
    });
    
    if (!response.ok) {
      throw new Error('Не удалось получить сообщение');
    }
    
    return await response.json();
  } catch (error) {
    console.error(`Error fetching message ${id}:`, error);
    return null;
  }
};

export const sendMessage = async (messageData) => {
  try {
    const response = await fetch(`${API_URL}/messages`, {
      method: 'POST',
      headers: getAuthHeader(),
      body: JSON.stringify(messageData)
    });
    
    if (!response.ok) {
      throw new Error('Не удалось отправить сообщение');
    }
    
    return await response.json();
  } catch (error) {
    console.error('Error sending message:', error);
    throw error;
  }
}; 
