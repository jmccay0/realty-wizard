import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

interface Hotspot {
  id: string;
  x: number; // percentage from left
  y: number; // percentage from top
  width: number; // percentage width
  height: number; // percentage height
  title: string;
  description: string;
}

// Hotspots for Page 0 (first page of TREC contract)
const page0Hotspots: Hotspot[] = [
  {
    id: 'parties',
    x: 5,
    y: 8,
    width: 90,
    height: 2.5,
    title: 'Parties to the Contract',
    description: 'Lists the Seller(s) and Buyer(s). Make sure all owners are listed as sellers, and all people buying should be listed as buyers. If married, both spouses typically need to be on the contract.',
  },
  {
    id: 'property-address',
    x: 5,
    y: 13,
    width: 90,
    height: 3,
    title: 'Property',
    description: 'The complete street address, city, county, and legal description of the property being sold. This must match the title exactly.',
  },
  {
    id: 'sales-price',
    x: 5,
    y: 30,
    width: 90,
    height: 8,
    title: '3. Sales Price',
    description: 'The total purchase price broken down into cash, financing, assumption of loan, and seller financing. All these components add up to the total sales price.',
  },
];

const TrecExplainer: React.FC = () => {
  const navigate = useNavigate();
  const [currentPage, setCurrentPage] = useState(0);
  const [activeHotspot, setActiveHotspot] = useState<string | null>(null);

  const totalPages = 11;

  const handleHotspotClick = (hotspot: Hotspot, event: React.MouseEvent) => {
    event.stopPropagation();
    setActiveHotspot(activeHotspot === hotspot.id ? null : hotspot.id);
  };

  const closeTooltip = () => {
    setActiveHotspot(null);
  };

  const getHotspotsForPage = (page: number): Hotspot[] => {
    const hotspotsByPage: Record<number, Hotspot[]> = {
      0: page0Hotspots,
      1: [
        {
          id: 'financing',
          x: 5,
          y: 5,
          width: 90,
          height: 15,
          title: '4. Financing',
          description: 'How you\'re paying for the property. Includes details about third-party financing, FHA/VA loans, and what happens if financing falls through. This section is critical if you need a loan.',
        },
        {
          id: 'earnest-money',
          x: 5,
          y: 45,
          width: 90,
          height: 10,
          title: '5. Earnest Money',
          description: 'Details about the earnest money deposit: how much, when it\'s due, who holds it, and what happens to it. This money shows you\'re serious and is applied to your down payment at closing.',
        },
      ],
      2: [
        {
          id: 'title-policy',
          x: 5,
          y: 5,
          width: 90,
          height: 12,
          title: '6. Title Policy and Survey',
          description: 'Seller agrees to provide title insurance and a survey. Title insurance protects you from ownership disputes. Survey shows property boundaries and any encroachments.',
        },
        {
          id: 'property-condition',
          x: 5,
          y: 30,
          width: 90,
          height: 15,
          title: '7. Property Condition',
          description: 'Explains that property is sold "as is" unless specifically excluded. Covers inspections, repairs, lender requirements, and who pays for what.',
        },
      ],
      3: [
        {
          id: 'brokers',
          x: 5,
          y: 5,
          width: 90,
          height: 15,
          title: '8. Brokers\' Fees',
          description: 'Details about real estate agent commissions - who pays them and how much. Typically seller pays both agents\' fees, but this can be negotiated.',
        },
        {
          id: 'closing',
          x: 5,
          y: 35,
          width: 90,
          height: 12,
          title: '9. Closing',
          description: 'When and where closing will occur, what costs each party pays, and how property taxes and HOA fees are split (prorated) between buyer and seller.',
        },
      ],
      4: [
        {
          id: 'possession',
          x: 5,
          y: 5,
          width: 90,
          height: 8,
          title: '10. Possession',
          description: 'When you get the keys! Typically at closing and funding, but can be earlier/later by agreement. Important: you don\'t own it until closing, even if you have keys.',
        },
        {
          id: 'special-provisions',
          x: 5,
          y: 20,
          width: 90,
          height: 10,
          title: '11. Special Provisions',
          description: 'Additional terms and conditions specific to this sale. This is where unique agreements go. Must not contradict other parts of the contract.',
        },
      ],
      5: [
        {
          id: 'settlement-disputes',
          x: 5,
          y: 5,
          width: 90,
          height: 15,
          title: '12. Settlement and Other Expenses',
          description: 'Who pays for what: title policy, escrow fees, tax certificates, HOA documents, etc. Some costs are customary, others are negotiable.',
        },
      ],
      6: [
        {
          id: 'prorations',
          x: 5,
          y: 5,
          width: 90,
          height: 12,
          title: '13. Prorations and Rollback Taxes',
          description: 'How property taxes, HOA dues, rents, and other costs are divided between buyer and seller based on the closing date.',
        },
        {
          id: 'casualty-loss',
          x: 5,
          y: 30,
          width: 90,
          height: 10,
          title: '14. Casualty Loss',
          description: 'What happens if the property is damaged or destroyed before closing. Usually seller bears the risk until closing.',
        },
      ],
      7: [
        {
          id: 'default',
          x: 5,
          y: 5,
          width: 90,
          height: 15,
          title: '15. Default',
          description: 'What happens if either party doesn\'t fulfill their obligations. Covers remedies, attorney fees, and earnest money disposition.',
        },
        {
          id: 'mediation',
          x: 5,
          y: 35,
          width: 90,
          height: 8,
          title: '16. Mediation',
          description: 'Agreement to try mediation before going to court if disputes arise. Mediation is cheaper and faster than litigation.',
        },
      ],
      8: [
        {
          id: 'attorney-fees',
          x: 5,
          y: 5,
          width: 90,
          height: 8,
          title: '17. Attorney\'s Fees',
          description: 'If there\'s a lawsuit, the losing party typically pays the winning party\'s attorney fees. This encourages settling disputes.',
        },
        {
          id: 'representations',
          x: 5,
          y: 25,
          width: 90,
          height: 12,
          title: '19. Representations',
          description: 'Both parties confirm they have authority to enter this contract and haven\'t filed bankruptcy. Important legal protections.',
        },
      ],
      9: [
        {
          id: 'notices',
          x: 5,
          y: 5,
          width: 90,
          height: 10,
          title: '21. Notices',
          description: 'How and where to send official communications. Email addresses, phone numbers, and addresses for all parties. Keep these current!',
        },
        {
          id: 'agreement-parties',
          x: 5,
          y: 30,
          width: 90,
          height: 12,
          title: '22. Agreement of Parties',
          description: 'Final provisions: entire agreement, modifications must be in writing, contract binding on heirs/successors, and other legal terms.',
        },
      ],
      10: [
        {
          id: 'consult-attorney',
          x: 5,
          y: 5,
          width: 90,
          height: 8,
          title: 'Consult an Attorney',
          description: 'Important advice: This is a legally binding contract. If you don\'t understand it, consult a real estate attorney before signing!',
        },
        {
          id: 'signatures',
          x: 5,
          y: 30,
          width: 90,
          height: 25,
          title: 'Signatures',
          description: 'Where all parties sign and date. The contract isn\'t valid until signed by all parties. The "effective date" starts when the last party signs and all parties are notified.',
        },
      ],
    };

    return hotspotsByPage[page] || [];
  };

  const activeHotspotData = activeHotspot
    ? getHotspotsForPage(currentPage).find(h => h.id === activeHotspot)
    : null;

  return (
    <div style={{
      minHeight: '100vh',
      backgroundColor: 'var(--realwiz-gray-50)',
      padding: '20px',
    }}>
      {/* Header */}
      <div style={{
        maxWidth: '1200px',
        margin: '0 auto 20px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
      }}>
        <button
          onClick={() => navigate('/')}
          className="btn-secondary"
        >
          ← Back to Home
        </button>
        <h1 style={{ fontSize: '24px', margin: 0 }}>
          Understanding Your Texas Real Estate Contract
        </h1>
      </div>

      {/* Page Navigation */}
      <div style={{
        maxWidth: '1200px',
        margin: '0 auto 20px',
        display: 'flex',
        gap: '10px',
        alignItems: 'center',
        justifyContent: 'center',
      }}>
        <button
          onClick={() => setCurrentPage(Math.max(0, currentPage - 1))}
          disabled={currentPage === 0}
          className="btn-secondary"
        >
          Previous
        </button>
        <span style={{ fontWeight: 'bold' }}>
          Page {currentPage + 1} of {totalPages}
        </span>
        <button
          onClick={() => setCurrentPage(Math.min(totalPages - 1, currentPage + 1))}
          disabled={currentPage === totalPages - 1}
          className="btn-secondary"
        >
          Next
        </button>
      </div>

      {/* Interactive Form Container */}
      <div style={{
        maxWidth: '1200px',
        margin: '0 auto',
        position: 'relative',
        backgroundColor: 'white',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)',
        borderRadius: '8px',
        overflow: 'hidden',
      }}>
        {/* Form Image */}
        <div style={{ position: 'relative', width: '100%' }}>
          <img
            src={`/trec-forms/trec-20-18-page-${currentPage}.png`}
            alt={`TREC Contract Page ${currentPage + 1}`}
            style={{
              width: '100%',
              height: 'auto',
              display: 'block',
            }}
          />

          {/* Clickable Hotspots */}
          {getHotspotsForPage(currentPage).map((hotspot) => (
            <div
              key={hotspot.id}
              onClick={(e) => handleHotspotClick(hotspot, e)}
              style={{
                position: 'absolute',
                left: `${hotspot.x}%`,
                top: `${hotspot.y}%`,
                width: `${hotspot.width}%`,
                height: `${hotspot.height}%`,
                backgroundColor: activeHotspot === hotspot.id
                  ? 'rgba(46, 125, 50, 0.3)'
                  : 'rgba(46, 125, 50, 0.15)',
                border: '2px solid var(--realwiz-green)',
                borderRadius: '4px',
                cursor: 'pointer',
                transition: 'all 0.2s',
                zIndex: activeHotspot === hotspot.id ? 10 : 1,
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.backgroundColor = 'rgba(46, 125, 50, 0.3)';
              }}
              onMouseLeave={(e) => {
                if (activeHotspot !== hotspot.id) {
                  e.currentTarget.style.backgroundColor = 'rgba(46, 125, 50, 0.15)';
                }
              }}
            >
              {/* Question mark indicator */}
              <div style={{
                position: 'absolute',
                top: '-10px',
                right: '-10px',
                width: '24px',
                height: '24px',
                borderRadius: '50%',
                backgroundColor: 'var(--realwiz-green)',
                color: 'white',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                fontSize: '14px',
                fontWeight: 'bold',
                boxShadow: '0 2px 4px rgba(0,0,0,0.2)',
              }}>
                ?
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Tooltip Modal */}
      {activeHotspotData && (
        <div
          onClick={closeTooltip}
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: 'rgba(0,0,0,0.5)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            zIndex: 1000,
            padding: '20px',
          }}
        >
          <div
            onClick={(e) => e.stopPropagation()}
            style={{
              backgroundColor: 'white',
              borderRadius: '8px',
              padding: '24px',
              maxWidth: '500px',
              width: '100%',
              boxShadow: '0 4px 16px rgba(0,0,0,0.2)',
            }}
          >
            <div style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'start',
              marginBottom: '16px',
            }}>
              <h3 style={{
                fontSize: '20px',
                fontWeight: 'bold',
                color: 'var(--realwiz-green)',
                margin: 0,
              }}>
                {activeHotspotData.title}
              </h3>
              <button
                onClick={closeTooltip}
                style={{
                  background: 'none',
                  border: 'none',
                  fontSize: '24px',
                  cursor: 'pointer',
                  color: 'var(--realwiz-gray-600)',
                  padding: 0,
                  width: '30px',
                  height: '30px',
                }}
              >
                ×
              </button>
            </div>
            <p style={{
              fontSize: '16px',
              lineHeight: '1.6',
              color: 'var(--realwiz-gray-700)',
              margin: 0,
            }}>
              {activeHotspotData.description}
            </p>
          </div>
        </div>
      )}

      {/* Help Text */}
      <div style={{
        maxWidth: '1200px',
        margin: '20px auto 0',
        textAlign: 'center',
        color: 'var(--realwiz-gray-600)',
        fontSize: '14px',
      }}>
        💡 Click on the highlighted areas to learn more about each section
      </div>
    </div>
  );
};

export default TrecExplainer;
