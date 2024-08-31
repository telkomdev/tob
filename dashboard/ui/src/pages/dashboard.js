import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

function Dashboard() {
  const [services, setServices] = useState([]);
  const [dashboardTitle, setDashboardTitle] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedTag, setSelectedTag] = useState('');
  const [shouldRedirect, setShouldRedirect] = useState(false);
  const [username, setUsername] = useState(localStorage.getItem('username'));
  const [token, setToken] = useState(localStorage.getItem('token'));
  const navigate = useNavigate();

  useEffect(() => {
    if (username == null) {
      setShouldRedirect(true);
    }
  }, [username]);

  useEffect(() => {
    const fetchServiceData = async () => {
      try {
        const response = await fetch('/api/services', {
            method: 'GET',
            headers: {
                'Authorization': token
            },
        });
        if (!response.ok) {
          setError(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();
        if (result.success) {
          const serviceArray = Object.keys(result.data.data).map(key => {
            const service = result.data.data[key];
            const latestCheckTime = new Date(Date.now() - service.checkInterval * 1000);

            if (service.tags) {
              service.tags.push(service.kind);
            }

            return {
              name: key,
              ...service,
              latestCheckTime: latestCheckTime.toLocaleString(),
            };
          });

          serviceArray.sort((a, b) => {
            if (a.status === 'DOWN' && b.status === 'UP') return -1;
            if (a.status === 'UP' && b.status === 'DOWN') return 1;
            return 0;
          });

          setServices(serviceArray);
          setDashboardTitle(result.data.dashboardTitle);
        } else {
          setError('Failed to retrieve services');
        }
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchServiceData();
    const intervalId = setInterval(fetchServiceData, 5000);
    return () => clearInterval(intervalId);
  }, [token]);

  const logout = () => {
    localStorage.removeItem('username');
    localStorage.removeItem('token');
    setShouldRedirect(true);
  };

  if (shouldRedirect) {
    return navigate('/');
  }

  const keyframes = `
  @keyframes pulse {
    0% {
      opacity: 1;
    }
    50% {
      opacity: 0.7;
    }
    100% {
      opacity: 1;
    }
  }
  `;

  const getStatusStyle = (status) => ({
    padding: '5px 10px',
    borderRadius: '4px',
    color: '#fff',
    fontWeight: 'bold',
    backgroundColor: status === 'UP' ? '#28a745' : '#dc3545',
    whiteSpace: 'nowrap',
    animation: 'pulse 1s infinite',
  });

  const getTagsStyle = () => ({
    display: 'flex',
    flexWrap: 'wrap',
    marginTop: '5px',
    gap: '5px',
  });

  const getTagStyle = () => ({
    backgroundColor: '#fff',
    color: '#333',
    borderRadius: '12px',
    padding: '3px 8px',
    fontSize: '12px',
    border: '1px solid #ddd',
    boxShadow: '0 1px 3px rgba(0, 0, 0, 0.1)',
    cursor: 'pointer',
  });

  const handleSearchTermChange = (event) => {
    setSearchTerm(event.target.value);
  };

  const handleTagClick = (tag) => {
    setSelectedTag(tag === selectedTag ? '' : tag);
  };

  const filteredServices = services.filter(service => {
    const matchesSearchTerm = service.name.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesTag = selectedTag ? service.tags.includes(selectedTag) : true;
    return matchesSearchTerm && matchesTag;
  });

  return (
    <div style={{ padding: '20px', fontFamily: 'Arial, sans-serif', backgroundColor: '#f4f4f4' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h1 style={{ textAlign: 'center', flex: 1 }}>{dashboardTitle}</h1>
        <button
          onClick={logout}
          style={{
            padding: '10px 20px',
            backgroundColor: '#007bff',
            color: '#fff',
            border: 'none',
            borderRadius: '5px',
            fontWeight: 'bold',
            cursor: 'pointer',
            boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
          }}
        >
          Logout
        </button>
      </div>

      <input
        type="text"
        placeholder="Search services..."
        value={searchTerm}
        onChange={handleSearchTermChange}
        style={{
          width: '80%',
          maxWidth: '600px',
          padding: '10px',
          marginBottom: '20px',
          borderRadius: '8px',
          border: '1px solid #ddd',
          boxShadow: '0 1px 3px rgba(0, 0, 0, 0.1)',
          display: 'block',
          marginLeft: 'auto',
          marginRight: 'auto',
        }}
      />
      {loading && <p>Loading services...</p>}
      {error && <p>Error: {error}</p>}
      {!loading && !error && (
        <ul style={{ maxWidth: '800px', margin: '0 auto', padding: '0', listStyle: 'none' }}>
          {filteredServices.map((service, index) => (
            <li key={index} style={{
              backgroundColor: '#fff',
              margin: '10px 0',
              padding: '15px',
              borderRadius: '8px',
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'flex-start',
              boxShadow: '0 2px 4px rgba(0, 0, 0, 0.1)',
            }}>
              
              <style>{keyframes}</style>
              <div style={{
                display: 'flex',
                justifyContent: 'space-between',
                width: '100%',
                marginBottom: '10px',
              }}>
                <span style={{ fontSize: '18px', fontWeight: 'bold' }}>{service.name}</span>
                <span style={getStatusStyle(service.status)}>
                  {service.status === 'UP' ? 'OK' : 'Not OK'}
                </span>
              </div>

              {service.status === 'DOWN' && (
                <span style={{ fontSize: '14px', color: '#dc3545', marginBottom: '10px' }}>
                  {service.messageDetails}
                </span>
              )}
              
              <span style={{ fontSize: '13px', color: '#555' }}>
                <span style={{ color: '#593f03' }}>Last checked: {service.latestCheckTime}</span> 
              </span>
              {service.tags && (
                <div style={getTagsStyle()}>
                  {service.tags.map((tag, tagIndex) => (
                    <span 
                      key={tagIndex} 
                      style={{
                        ...getTagStyle(),
                        backgroundColor: selectedTag === tag ? '#007bff' : '#fff',
                        color: selectedTag === tag ? '#fff' : '#333',
                      }}
                      onClick={() => handleTagClick(tag)}>
                      {tag}
                    </span>
                  ))}
                </div>
              )}
            </li>
          ))}
        </ul>
      )}

      <footer style={{
        position: 'fixed',
        bottom: '10px',
        right: '10px',
        fontSize: '15px',
        color: '#666',
      }}>
        Status Page by <a href="https://github.com/telkomdev/tob" target="_blank" rel="noopener noreferrer" style={{ color: '#007bff', textDecoration: 'none' }}>Tob</a>
      </footer>
    </div>
  );
}

export default Dashboard;