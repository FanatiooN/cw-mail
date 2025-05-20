import React, { useEffect, useState } from 'react';
import MailLayout from './layout/MailLayout';
import MessageList from './MessageList';
import { 
  getInboxMessages, 
  getSentMessages, 
  getSpamMessages,
  getTrashMessages 
} from '../services/api';

function MailPage({ type, title }) {
  const [messages, setMessages] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const loadMessages = async () => {
    try {
      setLoading(true);
      setError(null);
      
      let data = [];
      
      switch (type) {
        case 'inbox':
          data = await getInboxMessages();
          break;
        case 'sent':
          data = await getSentMessages();
          break;
        case 'spam':
          data = await getSpamMessages();
          break;
        case 'trash':
          data = await getTrashMessages();
          break;
        default:
          setError('Неизвестный тип папки');
      }
      
      setMessages(data);
    } catch (err) {
      console.error('Error loading messages:', err);
      setError('Не удалось загрузить сообщения');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadMessages();
  }, [type]);

  return (
    <MailLayout title={title}>
      <div className="mb-4 flex justify-between items-center">
        <h2 className="text-xl font-semibold">{title}</h2>
        <button 
          onClick={loadMessages} 
          className="px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 rounded-md"
        >
          Обновить
        </button>
      </div>
      
      <MessageList 
        messages={messages} 
        loading={loading} 
        error={error} 
      />
    </MailLayout>
  );
}

export default MailPage; 