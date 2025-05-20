import React from 'react';
import { useNavigate } from 'react-router-dom';
import Sidebar from './Sidebar';

function MailLayout({ children, title }) {
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/login");
  };

  return (
    <div className="flex h-screen w-full flex-col">
      <header className="flex h-14 items-center justify-between border-b bg-white px-4">
        <h1 className="text-xl font-bold md:hidden">CW Mail</h1>
        <div className="flex-1 md:ml-4">
          <h2 className="text-xl font-bold">{title}</h2>
        </div>
        <button
          onClick={handleLogout}
          className="ml-auto rounded-md px-3 py-1.5 text-sm font-medium text-gray-700 hover:bg-gray-100"
        >
          Выйти
        </button>
      </header>
      <div className="flex flex-1 overflow-hidden">
        <Sidebar />
        <main className="flex-1 overflow-y-auto bg-white">
          <div className="container mx-auto p-4 md:p-6">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
}

export default MailLayout; 