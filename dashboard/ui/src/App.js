import React, { useEffect, useState } from 'react';

function App() {
  const [services, setServices] = useState([]);
  const [dashboardTitle, setDashboardTitle] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedTag, setSelectedTag] = useState('');

  useEffect(() => {
    const fetchServiceData = async () => {
      try {
        const response = await fetch('/api/services');
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const result = await response.json();
        if (result.success) {
          const serviceArray = Object.keys(result.data.data).map(key => {
            const service = result.data.data[key];
            // Calculate the latest check time
            const latestCheckTime = new Date(Date.now() - service.checkInterval * 1000);

            // add service.kind to tags
            if (service.tags) {
              service.tags.push(service.kind);
            }

            return {
              name: key,
              ...service,
              latestCheckTime: latestCheckTime.toLocaleString(), // Convert to readable string
            };
          });

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

    fetchServiceData();
    const intervalId = setInterval(fetchServiceData, 5000);
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
      cursor: 'pointer', // To indicate it's clickable
    };
  };

  // Handle search term change
  const handleSearchTermChange = (event) => {
    setSearchTerm(event.target.value);
  };

  // Handle tag click to filter by tag
  const handleTagClick = (tag) => {
    setSelectedTag(tag === selectedTag ? '' : tag); // Toggle tag filter
  };

  // Filter services based on search term and selected tag
  const filteredServices = services.filter(service => {
    const matchesSearchTerm = service.name.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesTag = selectedTag ? service.tags.includes(selectedTag) : true;
    return matchesSearchTerm && matchesTag;
  });

  return (
    <div style={{ padding: '20px', fontFamily: 'Arial, sans-serif', backgroundColor: '#f4f4f4' }}>
      <h1 style={{ textAlign: 'center' }}>{dashboardTitle}</h1>
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
        <ul style={{ maxWidth: '600px', margin: '0 auto', padding: '0', listStyle: 'none' }}>
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
              <span style={{ fontSize: '12px', color: '#555' }}>
                Last checked: {service.latestCheckTime}
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
    </div>
  );
}

export default App;
