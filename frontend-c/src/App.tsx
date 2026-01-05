import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "@/contexts/AuthContext";
import { ToastProvider } from "@/contexts/ToastContext";
import Login from "@/pages/Login";
import FacultyDashboard from "@/pages/faculty/FacultyDashboard";
import CertificateVerification from "@/pages/faculty/CertificateVerification";
import FacultyProfile from "@/pages/faculty/FacultyProfile";
import HodDashboard from "@/pages/hod/HodDashboard";
import NotFound from "./pages/NotFound";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <AuthProvider>
      <ToastProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<Navigate to="/login" replace />} />
            <Route path="/login" element={<Login />} />
            <Route path="/faculty/dashboard" element={<FacultyDashboard />} />
            <Route path="/faculty/certificates" element={<CertificateVerification />} />
            <Route path="/faculty/profile" element={<FacultyProfile />} />
            <Route path="/hod/dashboard" element={<HodDashboard />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </BrowserRouter>
      </ToastProvider>
    </AuthProvider>
  </QueryClientProvider>
);

export default App;
