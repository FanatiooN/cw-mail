import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { ArrowLeft, Trash2 } from 'lucide-react';
import MailLayout from './layout/MailLayout';
import { sendMessage } from '../services/api';

function ComposeForm() {
  const [recipient, setRecipient] = useState("");
  const [subject, setSubject] = useState("");
  const [content, setContent] = useState("");
  const [isSelfDestruct, setIsSelfDestruct] = useState(false);
  const [readLimit, setReadLimit] = useState("1");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError("");
    
    try {
      await sendMessage({
        receiver_email: recipient,
        subject,
        body: content,
        read_limit: isSelfDestruct ? parseInt(readLimit) : 0
      });

      alert(isSelfDestruct 
        ? "Самоуничтожающееся сообщение успешно отправлено"
        : "Сообщение успешно отправлено");
      
      navigate("/sent");
    } catch (error) {
      console.error("Error sending message:", error);
      setError("Не удалось отправить сообщение");
    } finally {
      setIsLoading(false);
    }
  };

  const handleCancel = () => {
    navigate(-1);
  };

  return (
    <MailLayout title="Создать сообщение">
      <div className="w-full max-w-2xl mx-auto bg-white rounded-lg shadow">
        <div className="flex items-center justify-between p-4 border-b">
          <div className="flex items-center gap-2">
            <button onClick={handleCancel} className="p-2 hover:bg-gray-100 rounded-full">
              <ArrowLeft className="h-5 w-5" />
            </button>
            <h1 className="text-lg font-semibold">Новое сообщение</h1>
          </div>
          <button className="p-2 hover:bg-gray-100 rounded-full">
            <Trash2 className="h-5 w-5" />
          </button>
        </div>
        
        {error && (
          <div className="p-3 mx-4 mt-4 text-sm text-red-600 bg-red-50 rounded-md">
            {error}
          </div>
        )}
        
        <form onSubmit={handleSubmit} className="p-4 space-y-4">
          <div className="space-y-2">
            <label htmlFor="recipient" className="text-sm font-medium">
              Получатель
            </label>
            <input
              id="recipient"
              type="email"
              value={recipient}
              onChange={(e) => setRecipient(e.target.value)}
              placeholder="email@example.com"
              className="w-full p-2 border rounded-md"
              required
            />
          </div>
          
          <div className="space-y-2">
            <label htmlFor="subject" className="text-sm font-medium">
              Тема
            </label>
            <input
              id="subject"
              type="text"
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
              placeholder="Введите тему сообщения"
              className="w-full p-2 border rounded-md"
              required
            />
          </div>
          
          <div className="space-y-2">
            <label htmlFor="content" className="text-sm font-medium">
              Содержание
            </label>
            <textarea
              id="content"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="Введите текст сообщения"
              rows={10}
              className="w-full p-2 border rounded-md resize-none"
              required
            />
          </div>
          
          <div className="flex items-center space-x-2">
            <input
              type="checkbox"
              id="selfDestruct"
              checked={isSelfDestruct}
              onChange={(e) => setIsSelfDestruct(e.target.checked)}
              className="rounded border-gray-300"
            />
            <label htmlFor="selfDestruct" className="text-sm font-medium">
              Самоуничтожающееся сообщение
            </label>
          </div>
          
          {isSelfDestruct && (
            <div className="space-y-2">
              <label htmlFor="readLimit" className="text-sm font-medium">
                Лимит прочтений
              </label>
              <select
                value={readLimit}
                onChange={(e) => setReadLimit(e.target.value)}
                className="w-full p-2 border rounded-md"
              >
                <option value="1">1 прочтение</option>
                <option value="2">2 прочтения</option>
                <option value="3">3 прочтения</option>
                <option value="5">5 прочтений</option>
                <option value="10">10 прочтений</option>
              </select>
              <p className="text-xs text-gray-500">
                Сообщение будет автоматически удалено после {readLimit} прочтений или через 24 часа
              </p>
            </div>
          )}
          
          <div className="flex justify-end space-x-2 pt-4">
            <button
              type="button"
              onClick={handleCancel}
              className="px-4 py-2 border rounded-md hover:bg-gray-50"
            >
              Отмена
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
            >
              {isLoading ? "Отправка..." : "Отправить"}
            </button>
          </div>
        </form>
      </div>
    </MailLayout>
  );
}

export default ComposeForm; 