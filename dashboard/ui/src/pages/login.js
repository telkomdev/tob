import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

export default function Login(props) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [shouldRedirect, setShouldRedirect] = useState(false);
  const [disabledLogout, setDisabledLogout] = useState(false);
  const [messageLogin, setMessageLogin] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    const storedUsername = localStorage.getItem('username');
    if (storedUsername == null) {
      setDisabledLogout(true);
    }
  }, []);

  const handleSubmit = (e) => {
    e.preventDefault();
    loginService();
  };

  const loginService = async () => {
    try {
      const response = await fetch('/api/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          username: username,
          password: password,
        }),
      });

      const data = await response.json();
      if (response.status !== 200) {
        setMessageLogin(data.message);
      } else {
        localStorage.setItem('username', data.data.username);
        localStorage.setItem('token', data.data.jwtString);

        setUsername('');
        setPassword('');
        setShouldRedirect(true);
      }
    } catch (err) {
      console.error(err);
      setMessageLogin('Please try again later');
      setUsername('');
      setPassword('');
    }
  };

  if (shouldRedirect) {
    return navigate("/dashboard");
  }

  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      minHeight: '100vh',
      fontFamily: 'Arial, sans-serif',
      backgroundColor: '#f4f4f4',
    }}>
      <header style={{ textAlign: 'center', marginBottom: '20px', paddingTop: '30px' }}>
        <h1>Tob Dashboard Login</h1>
      </header>
      <main style={{
        flexGrow: 1, // Make the main content area grow to fill available space
        display: 'flex',
        flexDirection: 'column',
        justifyContent: 'flex-start',
        alignItems: 'center',
        paddingTop: '20px', // Adjust as needed
      }}>
        <form onSubmit={handleSubmit} style={{ maxWidth: '400px', width: '100%', margin: '0 auto' }}>
          {messageLogin && (
            <div style={{ color: 'red', marginBottom: '10px', textAlign: 'center' }}>
              {messageLogin}
            </div>
          )}
          <div style={{ marginBottom: '15px' }}>
            <label htmlFor="username" style={{ display: 'block', marginBottom: '5px' }}>
              Username
            </label>
            <input
              type="text"
              id="username"
              name="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              style={{
                width: '100%',
                padding: '8px',
                borderRadius: '4px',
                border: '1px solid #ccc',
                boxSizing: 'border-box',
              }}
            />
          </div>
          <div style={{ marginBottom: '20px' }}>
            <label htmlFor="password" style={{ display: 'block', marginBottom: '5px' }}>
              Password
            </label>
            <input
              type="password"
              id="password"
              name="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              style={{
                width: '100%',
                padding: '8px',
                borderRadius: '4px',
                border: '1px solid #ccc',
                boxSizing: 'border-box',
              }}
            />
          </div>
          <button
            type="submit"
            style={{
              width: '100%',
              padding: '10px',
              backgroundColor: '#28a745',
              color: '#fff',
              border: 'none',
              borderRadius: '4px',
              cursor: 'pointer',
              fontWeight: 'bold',
            }}
          >
            Login
          </button>
        </form>
      </main>
      <footer style={{
        position: 'fixed',
        bottom: '10px',
        right: '20px',
        fontSize: '15px',
        color: '#666',
      }}>
        Status Page by <a href="https://github.com/telkomdev/tob" target="_blank" rel="noopener noreferrer" style={{ color: '#007bff', textDecoration: 'none' }}>Tob</a>
      </footer>
    </div>
  );
}
