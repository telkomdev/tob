import React from 'react';
import ReactDOM from 'react-dom/client';
import reportWebVitals from './reportWebVitals';
import { Dashboard, ErrorPage, Login } from './pages'
import {
  createBrowserRouter,
  RouterProvider,
  Route,
  Routes,
} from "react-router-dom";

// function Root() {
//   return (
//     <Routes>
//       <Route path="/" element={<Login />} />
//       <Route path="/dashboard/*" element={<Dashboard />} />
//       <Route path="/about/*" element={<div>Hello</div>} />
//     </Routes>
//   );
// }

const router = createBrowserRouter([
  {
    path: "/",
    element: <Login/>,
    errorElement: <ErrorPage/>
  },
  {
    path: "/dashboard",
    element: <Dashboard/>,
    errorElement: <ErrorPage/>
  },
]);

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
