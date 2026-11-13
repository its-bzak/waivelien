import { useEffect, useState } from 'react'

// Strong typing that exactly matches our Go struct
interface InvoiceData {
  project_name: string;
  gc_name: string;
  subcontractor_name: string;
  invoice_amount: number;
  through_date: string;
}

function App() {
  const [invoice, setInvoice] = useState<InvoiceData | null>(null);
  const [loading, setLoading] = useState(false);

  // Fetch the data from Go when the page loads
  useEffect(() => {
    fetch('http://localhost:8080/api/invoice')
      .then(res => res.json())
      .then(data => setInvoice(data))
      .catch(err => console.error("Failed to fetch:", err));
  }, []);

  const handleSign = async () => {
    setLoading(true);
    try {
      const res = await fetch('http://localhost:8080/api/sign', { method: 'POST' });
      if (res.ok) {
        alert("Success! Check your backend folder for the generated PDF.");
      }
    } catch (err) {
      console.error("Failed to sign:", err);
      alert("Something went wrong.");
    } finally {
      setLoading(false);
    }
  };

  if (!invoice) return <div style={{ padding: "2rem" }}>Loading data from Go server...</div>;

  return (
    <div style={{ maxWidth: "600px", margin: "40px auto", fontFamily: "sans-serif" }}>
      <h1>Pending Lien Waiver</h1>
      
      <div style={{ background: "#f4f4f5", padding: "20px", borderRadius: "8px", marginBottom: "20px" }}>
        <p><strong>Project:</strong> {invoice.project_name}</p>
        <p><strong>General Contractor:</strong> {invoice.gc_name}</p>
        <p><strong>Amount:</strong> ${invoice.invoice_amount.toLocaleString()}</p>
        <p><strong>Through Date:</strong> {invoice.through_date}</p>
      </div>

      <button 
        onClick={handleSign} 
        disabled={loading}
        style={{ padding: "12px 24px", background: "#000", color: "#fff", border: "none", borderRadius: "6px", cursor: "pointer", fontSize: "16px" }}
      >
        {loading ? "Processing..." : "Sign & Release Funds"}
      </button>
    </div>
  )
}

export default App