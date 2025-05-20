import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import { ArrowLeft, Star, Trash2, Clock } from 'lucide-react';
import MailLayout from './layout/MailLayout';
import { getMessage } from '../services/api';

function MessageView() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [message, setMessage] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const fetchMessage = async () => {
      try {
        setLoading(true);
        const data = await getMessage(id);
        setMessage(data);
      } catch (err) {
        console.error(`Error fetching message ${id}:`, err);
        setError('Не удалось загрузить сообщение');
      } finally {
        setLoading(false);
      }
    };

    fetchMessage();
  }, [id]);

  const handleBack = () => {
    navigate(-1);
  };

  if (loading) {
    return (
      <MailLayout title="Загрузка сообщения...">
        <div className="flex justify-center py-10">
          <div className="animate-spin rounded-full h-10 w-10 border-t-2 border-b-2 border-blue-500"></div>
        </div>
      </MailLayout>
    );
  }

  if (error || !message) {
    return (
      <MailLayout title="Ошибка">
        <div className="p-4 text-red-600 bg-red-50 rounded-md">
          <p>{error || 'Сообщение не найдено'}</p>
          <button 
            onClick={handleBack} 
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md"
          >
            Вернуться назад
          </button>
        </div>
      </MailLayout>
    );
  }

  const date = new Date(message.created_at || message.sent_at || Date.now());
  const formattedDate = format(date, "d MMMM yyyy 'в' HH:mm", { locale: ru });
  const isSelfDestruct = message.read_limit > 0;

  return (
    <MailLayout title={message.subject}>
      <div className="bg-white rounded-lg shadow-sm border">
        <div className="p-4 border-b flex justify-between items-center">
          <button 
            onClick={handleBack} 
            className="p-2 hover:bg-gray-100 rounded-full"
          >
            <ArrowLeft className="h-5 w-5" />
          </button>
          
          <div className="flex space-x-2">
            <button className="p-2 hover:bg-gray-100 rounded-full">
              <Star className={`h-5 w-5 ${message.is_starred ? 'text-yellow-400 fill-yellow-400' : 'text-gray-400'}`} />
            </button>
            <button className="p-2 hover:bg-gray-100 rounded-full">
              <Trash2 className="h-5 w-5 text-gray-400" />
            </button>
          </div>
        </div>
        
        <div className="p-6">
          <div className="mb-6">
            <h1 className="text-2xl font-bold mb-4">{message.subject}</h1>
            
            <div className="flex items-center justify-between mb-2">
              <div>
                <p className="font-semibold">{message.sender_name || message.sender_email}</p>
                <p className="text-sm text-gray-500">Кому: {message.receiver_email}</p>
              </div>
              
              <div className="text-sm text-gray-500 flex items-center">
                {isSelfDestruct && (
                  <div className="flex items-center mr-3 text-orange-500">
                    <Clock className="w-4 h-4 mr-1" />
                    <span>Самоуничтожающееся ({message.read_limit} прочтений)</span>
                  </div>
                )}
                {formattedDate}
              </div>
            </div>
          </div>
          
          <div className="border-t pt-4 whitespace-pre-wrap">
            {message.body}
          </div>
          
          {message.attachments && message.attachments.length > 0 && (
            <div className="mt-6 border-t pt-4">
              <h3 className="text-sm font-semibold mb-2">Вложения ({message.attachments.length})</h3>
              <div className="space-y-2">
                {message.attachments.map((attachment, index) => (
                  <div key={index} className="border rounded-md p-2 flex items-center">
                    <span className="flex-grow truncate">{attachment.name}</span>
                    <a 
                      href={attachment.url} 
                      download
                      className="ml-2 px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 rounded-md"
                    >
                      Скачать
                    </a>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </MailLayout>
  );
}

export default MessageView; 