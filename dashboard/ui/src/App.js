import React, { useEffect, useState } from 'react';

function App() {
  const [services, setServices] = useState([]);
  const [dashboardTitle, setDashboardTitle] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Define the fetch function
    const fetchServiceData = async () => {
      try {
        const response = await fetch('/api/services');
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();
        if (result.success) {
          const serviceArray = Object.keys(result.data.data).map(key => ({
            name: key,
            ...result.data.data[key],
          }));

          // Sort the services: DOWN first, then UP
          serviceArray.sort((a, b) => {
            if (a.status === 'DOWN' && b.status === 'UP') return -1;
            if (a.status === 'UP' && b.status === 'DOWN') return 1;
            return 0;
          });

          setServices(serviceArray);
          setDashboardTitle(result.data.dashboardTitle);
        } else {
          throw new Error('Failed to retrieve services');
        }
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    // Initial fetch on component mount
    fetchServiceData();

    // Set an interval to fetch the data every X milliseconds
    const intervalId = setInterval(fetchServiceData, 5000); // Adjust the interval time as needed

    // Cleanup the interval on component unmount
    return () => clearInterval(intervalId);
  }, []);

  const getStatusStyle = (status) => {
    return {
      padding: '5px 10px',
      borderRadius: '4px',
      color: '#fff',
      fontWeight: 'bold',
      backgroundColor: status === 'UP' ? '#28a745' : '#dc3545',
      whiteSpace: 'nowrap',
    };
  };

  const getTagsStyle = () => {
    return {
      display: 'flex',
      flexWrap: 'wrap',
      marginTop: '5px',
      gap: '5px',
    };
  };

  const getTagStyle = () => {
    return {
      backgroundColor: '#fff',
      color: '#333',
      borderRadius: '12px',
      padding: '3px 8px',
      fontSize: '12px',
      border: '1px solid #ddd',
      boxShadow: '0 1px 3px rgba(0, 0, 0, 0.1)',
    };
  };

  return (
    <div style={{ padding: '20px', fontFamily: 'Arial, sans-serif', backgroundColor: '#f4f4f4' }}>
      <h1 style={{ textAlign: 'center' }}>{dashboardTitle}</h1>
      {loading && <p>Loading services...</p>}
      {error && <p>Error: {error}</p>}
      {!loading && !error && (
        <ul style={{ maxWidth: '600px', margin: '0 auto', padding: '0', listStyle: 'none' }}>
          {services.map((service, index) => (
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
              <div style={{
                display: 'flex',
                justifyContent: 'space-between',
                width: '100%',
                marginBottom: '10px',
              }}>
                <span style={{ fontSize: '18px' }}>{service.name}</span>
                <span style={getStatusStyle(service.status)}>
                  {service.status === 'UP' ? 'OK' : 'Not OK'}
                </span>
              </div>
              {service.tags && (
                <div style={getTagsStyle()}>
                  {service.tags.map((tag, tagIndex) => (
                    <span key={tagIndex} style={getTagStyle()}>
                      {tag}
                    </span>
                  ))}
                </div>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default App;
