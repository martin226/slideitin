import { motion } from "framer-motion";
import { Clock, FileText, Brain, Zap, Users } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

interface UseCaseProps {
  icon: React.ReactNode;
  title: string;
  description: string;
  className?: string;
}

const UseCase = ({ icon, title, description, className }: UseCaseProps) => {
  return (
    <motion.div 
      className={cn(
        "flex flex-col items-center p-6 rounded-lg bg-white shadow-md border border-gray-100 transition-all hover:shadow-lg",
        className
      )}
      whileHover={{ y: -5 }}
      initial={{ opacity: 0, y: 20 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ duration: 0.5 }}
    >
      <div className="w-12 h-12 rounded-full bg-amber-100 flex items-center justify-center mb-4">
        {icon}
      </div>
      <h3 className="text-xl font-semibold mb-2 text-gray-800">{title}</h3>
      <p className="text-gray-600 text-center">{description}</p>
    </motion.div>
  );
};

export default function UseCases() {
  const fadeInVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: (i: number) => ({
      opacity: 1,
      y: 0,
      transition: {
        duration: 0.7,
        delay: 0.1 + i * 0.1,
        ease: [0.25, 0.4, 0.25, 1],
      },
    }),
  };

  return (
    <section className="w-full py-20 bg-amber-50">
      <div className="container mx-auto px-4 md:px-6">
        {/* Section Header */}
        <motion.div 
          className="text-center max-w-3xl mx-auto mb-16"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true }}
          variants={fadeInVariants}
          custom={0}
        >
          <h2 className="text-3xl md:text-4xl font-bold mb-4 text-gray-800">
            Transform Your Content Into Beautiful Presentations
          </h2>
          <p className="text-lg text-gray-600">
            Slide it In helps you create professional slides in minutes, not hours.
            Here&apos;s how people are using it:
          </p>
        </motion.div>

        {/* Use Cases Grid with flexbox for better centering */}
        <div className="flex flex-wrap justify-center gap-6 mb-16 px-4">
          {/* First row of use cases */}
          <div className="flex-[0_0_100%] sm:flex-[0_0_calc(50%-12px)] lg:flex-[0_0_calc(33.333%-16px)]">
            <UseCase 
              icon={<Brain className="w-6 h-6 text-amber-600" />}
              title="Study Notes to Slides"
              description="Turn your dense study materials into clear, organized slides that make learning and revision easier."
              className="md:transform md:rotate-[-1deg] border-amber-200 h-full"
            />
          </div>
          
          <div className="flex-[0_0_100%] sm:flex-[0_0_calc(50%-12px)] lg:flex-[0_0_calc(33.333%-16px)]">
            <UseCase 
              icon={<Clock className="w-6 h-6 text-amber-600" />}
              title="Quick Presentation Drafts"
              description="Create the first draft of your presentation in minutes, then fine-tune as needed."
              className="md:transform md:translate-y-4 border-amber-200 h-full"
            />
          </div>
          
          <div className="flex-[0_0_100%] sm:flex-[0_0_calc(50%-12px)] lg:flex-[0_0_calc(33.333%-16px)]">
            <UseCase 
              icon={<FileText className="w-6 h-6 text-amber-600" />}
              title="Research Papers to Talks"
              description="Convert academic papers and research findings into presentation-ready formats instantly."
              className="md:transform md:rotate-[1deg] border-amber-200 h-full"
            />
          </div>
          
          <div className="flex-[0_0_100%] sm:flex-[0_0_calc(50%-12px)] lg:flex-[0_0_calc(33.333%-16px)]">
            <UseCase 
              icon={<Users className="w-6 h-6 text-amber-600" />}
              title="Client Proposals"
              description="Transform your proposal documents into professional slides that impress potential clients."
              className="md:transform md:rotate-[1deg] border-amber-200 h-full"
            />
          </div>
          
          <div className="flex-[0_0_100%] sm:flex-[0_0_calc(50%-12px)] lg:flex-[0_0_calc(33.333%-16px)]">
            <UseCase 
              icon={<Zap className="w-6 h-6 text-amber-600" />}
              title="Meeting Prep in Minutes"
              description="Turn meeting agendas and notes into structured slides for better team communication."
              className="md:transform md:translate-y-4 border-amber-200 h-full"
            />
          </div>
        </div>
        
        {/* CTA Section */}
        <motion.div 
          className="text-center"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true }}
          variants={fadeInVariants}
          custom={6}
        >
          <Button 
            className="px-8 py-6 bg-amber-500 hover:bg-amber-600 text-white rounded-full text-lg font-medium"
            onClick={() => window.location.href = '/start'}
          >
            Make a Presentation! (It&apos;s free)
          </Button>
          <p className="mt-4 text-gray-600 text-sm">No sign-up required. Start creating in seconds.</p>
        </motion.div>
      </div>
    </section>
  );
} 