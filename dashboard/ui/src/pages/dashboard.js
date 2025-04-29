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

        // Reset error state before each fetch attempt
        setError(null);

        const response = await fetch('/api/services', {
            method: 'GET',
            headers: {
                'Authorization': token
            },
        });

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

          const statusPriority = {
            DOWN: 0,
            MONITORED: 1,
            UP: 2,
          };
          
          serviceArray.sort((a, b) => {
            return statusPriority[a.status] - statusPriority[b.status];
          });

          setServices(serviceArray);
          setDashboardTitle(result.data.dashboardTitle);
        } else {
          if (result.message) {
            // if token is expired, force logout
            if (result.message.includes('token expired')) {
              logout();
            } else {
              setError(result.message);
            }
          } else {
            setError('Failed to retrieve services');
          }
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
      opacity: 0.2;
    }
    100% {
      opacity: 1;
    }
  }
  `;

  const getStatusStyle = (status) => {
    const randomDuration = Math.random() * 0.2 + 0.7;

    return {
      padding: '5px 10px',
      borderRadius: '4px',
      color: '#000',
      fontWeight: 'bold',
      backgroundColor:
        status === 'UP' ? '#28a745' :
        status === 'DOWN' ? '#dc3545' :
        '#ffc107', 
      whiteSpace: 'nowrap',
      animation: `pulse ${randomDuration.toFixed(2)}s infinite`,
      flexShrink: 0,
    };
  };
  
  
  
  
  const getTagsStyle = () => ({
    display: 'flex',
    flexWrap: 'wrap',
    marginTop: '5px',
    gap: '5px',
  });

  const getTagStyle = () => ({
    backgroundColor: 'rgba(255, 255, 255, 0.08)',
    color: '#ccc',
    borderRadius: '8px',
    padding: '2px 6px',
    fontSize: '12px',
    border: '1px solid rgba(255, 255, 255, 0.15)',
    fontWeight: '500',
    cursor: 'pointer',
  });

  const getSeverityColor = (line) => {
    if (line.includes('Warning')) return '#967205';
    if (line.includes('Critical') || line.includes('Danger')) return '#dc3545';
    if (line.includes('Info')) return '#04a0bf';
    return '#212529';
  };

  const handleSearchTermChange = (event) => {
    setSearchTerm(event.target.value);
  };

  const handleTagClick = (tag) => {
    setSelectedTag(tag === selectedTag ? '' : tag);
  };

  const filteredServices = services.filter(service => {
    const matchesSearchTerm = service.name.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesTagSearch = service.tags.some(tag => tag.toLowerCase().includes(searchTerm.toLowerCase()));
    const matchesPicsSearch = service.pics.some(pic => pic.toLowerCase().includes(searchTerm.toLowerCase()));
    const matchesTag = selectedTag ? service.tags.includes(selectedTag) : true;
    return (matchesSearchTerm || matchesTagSearch || matchesPicsSearch) && matchesTag;
  });

  return (
    <div style={{ padding: '20px', fontFamily: 'Arial, sans-serif', backgroundColor: '#212121', color: '#f5f5f5', minHeight: '100vh' }}>
      <div style={{ 
        display: 'flex', 
        justifyContent: 'space-between', 
        alignItems: 'center', 
        flexWrap: 'wrap',
        marginBottom: '20px'
      }}>
        <div style={{ flex: '1 1 auto', marginBottom: '10px' }}>
          <img
            src="tob.png"
            title='Tob the Bot (https://github.com/telkomdev/tob)'
            alt="Logo"
            style={{
              width: '90px',
              height: 'auto',
              borderRadius: '50%',
            }}
          />
        </div>
        <div style={{ flex: '1 1 auto', textAlign: 'center' }}>
          <h2>{dashboardTitle}</h2>
        </div>
        <div style={{ flex: '1 1 auto', textAlign: 'right' }}>
          <button
            onClick={logout}
            style={{
              padding: '10px 20px',
              backgroundColor: '#04a0bf',
              color: '#fff',
              border: 'none',
              borderRadius: '5px',
              fontWeight: 'bold',
              cursor: 'pointer',
              boxShadow: '0 2px 4px rgba(0, 0, 0, 0.5)',
            }}
          >
            Logout
          </button>
        </div>
      </div>

      <div style={{ display: 'flex', justifyContent: 'center', marginBottom: '20px' }}>
        <input
          type="text"
          placeholder="Search services..."
          value={searchTerm}
          onChange={handleSearchTermChange}
          style={{
            width: '80%',
            maxWidth: '600px',
            padding: '10px',
            borderRadius: '8px',
            border: '1px solid #444',
            backgroundColor: '#1e1e1e',
            color: '#f5f5f5',
            boxShadow: '0 1px 3px rgba(0, 0, 0, 0.6)',
            margin: '0 auto',
          }}
        />
      </div>

      {loading && <p>Loading services...</p>}
      {error && <p>Error: {error}</p>}
      {!loading && !error && (
        <ul style={{ maxWidth: '800px', margin: '0 auto', padding: '0', listStyle: 'none' }}>
        {filteredServices.map((service, index) => (
          <li key={index} style={{
            backgroundColor: '#353535',
            margin: '10px 0',
            padding: '15px',
            borderRadius: '8px',
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'flex-start',
            boxShadow: '0 2px 8px rgba(0, 0, 0, 0.6)',
          }}>
            
            <style>{keyframes}</style>
            <div style={{
              display: 'flex',
              justifyContent: 'space-between',
              width: '100%',
              marginBottom: '10px',
              flexWrap: 'wrap',
            }}>
              <span style={{ 
                fontSize: '18px', 
                fontWeight: 'bold', 
                flexGrow: 1, 
                marginRight: '10px', 
                wordWrap: 'break-word', // Allow breaking long words
                overflowWrap: 'break-word', // IE and Edge support
                maxWidth: 'calc(100% - 120px)', // Adjust based on status box size
              }}>
                {service.name}
              </span>
              <span style={getStatusStyle(service.status)}>
                {service.status === 'UP'
                  ? 'OK'
                  : service.status === 'DOWN'
                  ? 'Not OK'
                  : service.status === 'MONITORED'
                  ? 'Monitored'
                  : service.status}
              </span>
            </div>

            {service.messageDetails && service.kind === 'sslstatus' ? (
              <div style={{ fontSize: '11px', marginBottom: '10px' }}>
                <table style={{ width: '100%', borderCollapse: 'collapse', border: '1px solid #ccc' }}>
                  <thead>
                    <tr>
                      <th style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'left' }}>Severity</th>
                      <th style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'left' }}>Domain</th>
                      <th style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'left' }}>Message</th>
                      <th style={{ border: '1px solid #ccc', padding: '8px', textAlign: 'left' }}>Expiration Date</th>
                    </tr>
                  </thead>
                  <tbody>
                    {service.messageDetails.split('\n').map((line, index) => {
                      const parts = line.split('|').map(part => part.trim());

                      if (parts.length >= 4) {
                        const status = parts[0];
                        const domain = parts[1];
                        const remainingTime = parts[2];
                        const detail = parts[3];

                        // color by status
                        let rowColor = '';
                        if (status.toLowerCase().includes('danger')) {
                          rowColor = '#dc3545'; // red
                        } else if (status.toLowerCase().includes('warning')) {
                          rowColor = '#ffc107'; // yellow
                        } else if (status.toLowerCase().includes('info')) {
                          rowColor = '#28a745'; // green
                        }

                        return (
                          <tr key={index} style={{ color: rowColor }}>
                            <td style={{ border: '1px solid #ccc', padding: '8px' }}>{status}</td>
                            <td style={{ border: '1px solid #ccc', padding: '8px' }}>{domain}</td>
                            <td style={{ border: '1px solid #ccc', padding: '8px' }}>{remainingTime}</td>
                            <td style={{ border: '1px solid #ccc', padding: '8px' }}>{detail}</td>
                          </tr>
                        );
                      } else if (parts.length === 3) {
                        // handle when there is no date or time
                        const status = parts[0];
                        const domain = parts[1];
                        const detail = parts[2];

                        let rowColor = '';
                        if (status.toLowerCase().includes('danger')) {
                          rowColor = '#dc3545';
                        } else if (status.toLowerCase().includes('warning')) {
                          rowColor = '#ffc107';
                        } else if (status.toLowerCase().includes('info')) {
                          rowColor = '#28a745';
                        }

                        return (
                          <tr key={index} style={{ color: rowColor }}>
                            <td style={{ border: '1px solid #ccc', padding: '8px' }}>{status}</td>
                            <td style={{ border: '1px solid #ccc', padding: '8px' }}>{domain}</td>
                            <td style={{ border: '1px solid #ccc', padding: '8px' }}></td> {}
                            <td style={{ border: '1px solid #ccc', padding: '8px' }}>{detail}</td>
                          </tr>
                        );
                      }
                      return null;
                    })}
                  </tbody>
                </table>
              </div>
            )
            :
            (
              <div style={{ 
                fontSize: '14px', 
                color: '#dc3545', 
                marginBottom: '10px', 
                wordWrap: 'break-word',
                overflowWrap: 'break-word',
                whiteSpace: 'pre-line',
              }}>
                {service.messageDetails}
              </div>
            )
          }

            {service.pics && service.pics.length > 0 && (
              <span style={{ fontSize: '14px', color: '#555', marginBottom: '10px', fontWeight: 'bold' }}>
                <span style={{ color: '#04a0bf' }}>PICs: {service.pics.join(', ')}</span> 
              </span>
            )}

            <span style={{ fontSize: '12px', color: '#aaa' }}>
              <span style={{ color: '#d4af37', fontWeight: 500 }}>
                Last checked:
              </span>{' '}
              {service.latestCheckTime}
            </span>
            
            {service.tags && (
              <div style={getTagsStyle()}>
                {service.tags.map((tag, tagIndex) => (
                  <span 
                    key={tagIndex} 
                    style={{
                      ...getTagStyle(),
                      backgroundColor: selectedTag === tag ? '#04a0bf' : '#fff',
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
        right: '20px',
        fontSize: '15px',
        color: '#666',
      }}>
        Status Page by <a href="https://github.com/telkomdev/tob" target="_blank" rel="noopener noreferrer" style={{ color: '#007bff', textDecoration: 'none' }}>Tob</a>
      </footer>
    </div>
  );
}

export default Dashboard;
