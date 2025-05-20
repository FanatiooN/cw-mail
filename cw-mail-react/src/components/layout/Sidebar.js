import React from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { Inbox, Send, AlertTriangle, Trash2, Mail } from 'lucide-react';

function Sidebar() {
  const navigate = useNavigate();
  const location = useLocation();

  const menuItems = [
    { 
      title: "Входящие", 
      path: "/inbox", 
      icon: <Inbox className="h-5 w-5" /> 
    },
    { 
      title: "Отправленные", 
      path: "/sent", 
      icon: <Send className="h-5 w-5" /> 
    },
    { 
      title: "Спам", 
      path: "/spam", 
      icon: <AlertTriangle className="h-5 w-5" /> 
    },
    { 
      title: "Удаленные", 
      path: "/trash", 
      icon: <Trash2 className="h-5 w-5" /> 
    }
  ];

  return (
    <div className="h-full w-64 border-r bg-gray-50">
      <div className="flex h-14 items-center border-b px-4">
        <Mail className="mr-2 h-6 w-6" />
        <span className="text-lg font-medium">CW Mail</span>
      </div>
      <nav className="space-y-1 p-2">
        {menuItems.map((item) => (
          <button
            key={item.path}
            onClick={() => navigate(item.path)}
            className={`flex w-full items-center space-x-3 rounded-md px-3 py-2 text-left text-sm font-medium
              ${location.pathname === item.path 
                ? "bg-gray-200 text-gray-900" 
                : "text-gray-700 hover:bg-gray-100"}
            `}
          >
            {item.icon}
            <span>{item.title}</span>
          </button>
        ))}
      </nav>
      <div className="absolute bottom-5 w-64 px-4">
        <button
          onClick={() => navigate("/compose")}
          className="flex w-full items-center justify-center rounded-md bg-blue-600 px-3 py-2 text-white hover:bg-blue-700"
        >
          Написать
        </button>
      </div>
    </div>
  );
}

export default Sidebar; 