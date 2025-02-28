import { useState, useEffect, useRef, useMemo } from 'react';
import { Document, Page, pdfjs } from 'react-pdf';
import 'react-pdf/dist/esm/Page/AnnotationLayer.css';
import 'react-pdf/dist/esm/Page/TextLayer.css';
import { ChevronLeft, ChevronRight, Download, Edit } from 'lucide-react';

// Set worker path to a local path that Next.js can serve
pdfjs.GlobalWorkerOptions.workerSrc = `//cdn.jsdelivr.net/npm/pdfjs-dist@${pdfjs.version}/build/pdf.worker.min.mjs`;

interface PDFViewerClientProps {
    fileUrl: string;
    onDownload?: () => void;
    onEdit?: () => void;
}

function toDataURL(url: string) {
    return fetch(url)
        .then(response => response.blob())
        .then(blob => new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onloadend = () => resolve(reader.result);
            reader.onerror = reject;
            reader.readAsDataURL(blob);
        }));
}

export function PDFViewerClient({ fileUrl, onDownload, onEdit }: PDFViewerClientProps) {
    const [numPages, setNumPages] = useState<number>(0);
    const [pageNumber, setPageNumber] = useState<number>(1);
    const [loading, setLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);
    const [dataUrl, setDataUrl] = useState<string | null>(null);
    const containerRef = useRef<HTMLDivElement>(null);
    const documentLoaded = useRef<boolean>(false); // Track if document is loaded
    
    // Memoize the PDF.js options to prevent unnecessary rerenders
    const pdfOptions = useMemo(() => ({
        cMapUrl: `https://cdn.jsdelivr.net/npm/pdfjs-dist@${pdfjs.version}/cmaps/`,
        cMapPacked: true,
    }), []);
    
    useEffect(() => {
        const loadPdf = async () => {
            try {
                setLoading(true);
                // Convert the fileUrl to a data URL to avoid CORS issues
                const pdfDataUrl = await toDataURL(fileUrl);
                setDataUrl(pdfDataUrl as string);
                setError(null);
            } catch (err) {
                console.error('Error loading PDF:', err);
                setError('Failed to load PDF. Please try again later.');
            } finally {
                setLoading(false);
            }
        };
        
        loadPdf();
        
        // Cleanup function to ensure any remaining worker tasks are cleanly terminated
        return () => {
            // Reset state on unmount to prevent state updates on unmounted component
            documentLoaded.current = false;
        };
    }, [fileUrl]);
    
    function onDocumentLoadSuccess({ numPages }: { numPages: number }) {
        setNumPages(numPages);
        setLoading(false);
        documentLoaded.current = true;
    }
    
    function changePage(offset: number) {
        const newPageNumber = pageNumber + offset;
        if (newPageNumber >= 1 && newPageNumber <= numPages) {
            setPageNumber(newPageNumber);
        }
    }
    
    function previousPage() {
        changePage(-1);
    }
    
    function nextPage() {
        changePage(1);
    }
    
    // Separate click handlers for left and right navigation
    function handleLeftClick() {
        if (documentLoaded.current && pageNumber > 1) {
            previousPage();
        }
    }
    
    function handleRightClick() {
        if (documentLoaded.current && pageNumber < numPages) {
            nextPage();
        }
    }
    
    // Handle errors that might occur during rendering
    const handleError = (error: Error) => {
        console.error('PDF rendering error:', error);
        // Only update state if the error is not related to worker termination or transport
        if (!error.message.includes('Worker task was terminated') && 
            !error.message.includes('Transport destroyed')) {
            setError('Error rendering PDF: ' + error.message);
        }
    };
    
    if (loading && !dataUrl) {
        return (
            <div className="flex items-center justify-center w-full h-[300px]">
                <div className="text-gray-500">Loading PDF...</div>
            </div>
        );
    }
    
    if (error) {
        return (
            <div className="flex items-center justify-center w-full h-[300px]">
                <div className="text-red-500">{error}</div>
            </div>
        );
    }
    
    // Calculate which pages to render (previous, current, and next)
    const pagesToRender = [
        pageNumber > 1 ? pageNumber - 1 : null, // Previous page or null if on first page
        pageNumber, // Current page
        pageNumber < numPages ? pageNumber + 1 : null, // Next page or null if on last page
    ].filter(Boolean) as number[];
    
    return (
        <div ref={containerRef} className="w-full">
            {/* PDF Viewer with dedicated click areas for navigation */}
            <div className="relative w-full">
                {dataUrl && (
                    <Document
                        file={dataUrl}
                        onLoadSuccess={onDocumentLoadSuccess}
                        onLoadError={handleError}
                        loading={
                            <div className="flex items-center justify-center w-full h-64">
                                <div className="text-gray-500">Loading PDF...</div>
                            </div>
                        }
                        error={
                            <div className="flex items-center justify-center w-full h-64">
                                <div className="text-red-500">Failed to load PDF</div>
                            </div>
                        }
                        className="flex justify-center"
                        options={pdfOptions}
                        externalLinkTarget="_blank"
                    >
                        <div className="pdf-container flex justify-center">
                            {pagesToRender.map((pageNum) => (
                                <div 
                                    key={`page-${pageNum}`} 
                                    style={{
                                        display: pageNum === pageNumber ? 'block' : 'none'
                                    }}
                                    className="w-full"
                                >
                                    <Page 
                                        pageNumber={pageNum} 
                                        renderTextLayer={false} 
                                        renderAnnotationLayer={false}
                                        width={containerRef.current?.clientWidth || undefined}
                                        scale={1}
                                        loading={
                                            <div className="flex items-center justify-center w-full h-[400px]">
                                                <div className="text-gray-500">Loading page...</div>
                                            </div>
                                        }
                                        error={<div>Error loading page</div>}
                                        onRenderError={handleError}
                                    />
                                </div>
                            ))}
                        </div>
                    </Document>
                )}
                
                {/* Navigation overlays with dedicated click handlers */}
                <button 
                    onClick={handleLeftClick}
                    disabled={pageNumber <= 1}
                    className="absolute left-0 top-0 h-full w-[30%] z-20 cursor-pointer disabled:cursor-default opacity-0"
                    aria-label="Previous page"
                />
                
                <button 
                    onClick={handleRightClick}
                    disabled={pageNumber >= numPages}
                    className="absolute right-0 top-0 h-full w-[30%] z-20 cursor-pointer disabled:cursor-default opacity-0" 
                    aria-label="Next page"
                />
            </div>
            
            {/* Controls and buttons - flex column on mobile, row on desktop */}
            <div className="mt-6 flex flex-col sm:flex-row items-center justify-center w-full gap-4">
                {/* Navigation controls */}
                <div className="flex items-center justify-center space-x-6">
                    <button
                        onClick={previousPage}
                        disabled={pageNumber <= 1}
                        className="p-2 rounded-full hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                        aria-label="Previous page"
                    >
                        <ChevronLeft size={18} />
                    </button>
                    <span className="text-sm text-gray-600">
                        {pageNumber} / {numPages}
                    </span>
                    <button
                        onClick={nextPage}
                        disabled={pageNumber >= numPages}
                        className="p-2 rounded-full hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
                        aria-label="Next page"
                    >
                        <ChevronRight size={18} />
                    </button>
                </div>
                
                {/* Action buttons */}
                <div className="flex flex-col sm:flex-row gap-3">
                    {onDownload && (
                        <button 
                            className="py-2 px-4 rounded-lg bg-amber-500 hover:bg-amber-600 transition-colors flex items-center justify-center gap-2 text-white font-medium"
                            onClick={onDownload}
                        >
                            <Download size={16} />
                            Download as PDF
                        </button>
                    )}
                    
                    {onEdit && (
                        <button 
                            onClick={onEdit}
                            className="py-2 px-4 rounded-lg bg-amber-400 hover:bg-amber-500 transition-colors flex items-center justify-center gap-2 text-white font-medium"
                        >
                            <Edit size={16} />
                            Edit in PowerPoint
                        </button>
                    )}
                </div>
            </div>
        </div>
    );
}