import React from 'react';
import { Link } from 'react-router-dom';
import { formatDistanceToNow } from 'date-fns';
import { ru } from 'date-fns/locale';
import { Star, Clock } from 'lucide-react';

function MessageItem({ message }) {
  const isUnread = !message.read;
  const date = new Date(message.created_at || message.sent_at || Date.now());
  const formattedDate = formatDistanceToNow(date, { addSuffix: true, locale: ru });
  const isSelfDestruct = message.read_limit > 0;

  return (
    <Link 
      to={`/messages/${message.id}`}
      className={`block border-b hover:bg-gray-50 ${isUnread ? 'bg-blue-50' : ''}`}
    >
      <div className="flex items-center p-3">
        <div className="flex-shrink-0 mr-3">
          {isUnread && (
            <div className="w-2 h-2 bg-blue-600 rounded-full" />
          )}
        </div>
        
        <div className="flex-grow min-w-0">
          <div className="flex justify-between items-center mb-1">
            <div className="font-semibold truncate">
              {message.sender_email || message.receiver_email}
            </div>
            <div className="text-xs text-gray-500 flex items-center">
              {isSelfDestruct && (
                <Clock className="inline-block w-3 h-3 mr-1 text-orange-500" />
              )}
              {formattedDate}
            </div>
          </div>
          
          <div className="flex justify-between">
            <p className="text-sm font-medium truncate">{message.subject}</p>
            {message.is_starred && (
              <Star className="w-4 h-4 text-yellow-400" />
            )}
          </div>
          
          <p className="text-xs text-gray-500 truncate">{message.body}</p>
        </div>
      </div>
    </Link>
  );
}

function MessageList({ messages, loading, error }) {
  if (loading) {
    return (
      <div className="flex justify-center py-10">
        <div className="animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-blue-500"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="p-4 text-red-600 bg-red-50 rounded-md">
        <p>Ошибка при загрузке сообщений: {error}</p>
      </div>
    );
  }

  if (!messages || messages.length === 0) {
    return (
      <div className="p-10 text-center text-gray-500">
        <p>В этой папке нет сообщений</p>
      </div>
    );
  }

  return (
    <div className="border rounded-md overflow-hidden divide-y">
      {messages.map((message) => (
        <MessageItem key={message.id} message={message} />
      ))}
    </div>
  );
}

export default MessageList; 