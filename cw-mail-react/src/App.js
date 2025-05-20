import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import ComposeForm from './components/ComposeForm';
import Login from './components/Login';
import Register from './components/Register';
import ProtectedRoute from './components/ProtectedRoute';
import MailPage from './components/MailPage';
import MessageView from './components/MessageView';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/inbox" replace />} />
        <Route 
          path="/compose" 
          element={
            <ProtectedRoute>
              <ComposeForm />
            </ProtectedRoute>
          } 
        />
        <Route 
          path="/messages/:id" 
          element={
            <ProtectedRoute>
              <MessageView />
            </ProtectedRoute>
          } 
        />
        <Route 
          path="/inbox" 
          element={
            <ProtectedRoute>
              <MailPage type="inbox" title="Входящие" />
            </ProtectedRoute>
          } 
        />
        <Route 
          path="/sent" 
          element={
            <ProtectedRoute>
              <MailPage type="sent" title="Отправленные" />
            </ProtectedRoute>
          } 
        />
        <Route 
          path="/spam" 
          element={
            <ProtectedRoute>
              <MailPage type="spam" title="Спам" />
            </ProtectedRoute>
          } 
        />
        <Route 
          path="/trash" 
          element={
            <ProtectedRoute>
              <MailPage type="trash" title="Корзина" />
            </ProtectedRoute>
          } 
        />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
      </Routes>
    </Router>
  );
}

export default App;
