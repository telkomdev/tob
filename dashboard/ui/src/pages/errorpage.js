import React from 'react';
import { Link } from 'react-router-dom';

export default function ErrorPage() {
  return (
    <div style={styles.container}>
      <div style={styles.content}>
        <h1 style={styles.title}>404</h1>
        <p style={styles.message}>Oops! The page you're looking for doesn't exist.</p>
        <Link to="/" style={styles.link}>Go back to the homepage</Link>
      </div>
    </div>
  );
}

const styles = {
  container: {
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '100vh',
    backgroundColor: '#f4f4f4',
    fontFamily: 'Arial, sans-serif',
  },
  content: {
    textAlign: 'center',
    maxWidth: '400px',
    width: '100%',
    backgroundColor: '#fff',
    padding: '40px',
    borderRadius: '8px',
    boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)',
  },
  title: {
    fontSize: '72px',
    marginBottom: '20px',
    color: '#333',
  },
  message: {
    fontSize: '18px',
    marginBottom: '20px',
    color: '#555',
  },
  link: {
    fontSize: '16px',
    color: '#007bff',
    textDecoration: 'none',
    border: '1px solid #007bff',
    padding: '10px 20px',
    borderRadius: '4px',
  },
};
